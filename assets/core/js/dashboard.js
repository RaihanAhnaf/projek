var statussummaries = [];

var summaryVM = kendo.observable({
    new: 0,
    inProgress: 0,
    validated: 0,
    resolution: 0,
    selected: "New Contract",
    addhoverclass: function(classname){
        $("#applicationStatus li.active").removeClass("active")
        $("#applicationStatus li." + classname).addClass("active")
    },
    newContractClicked: function(){
        summaryVM.set("selected", "New Contract")
        summaryVM.addhoverclass("classcontract")
        refreshContractGrid()
    },
    progressClicked: function(){
        summaryVM.set("selected", "Progress")
        summaryVM.addhoverclass("classprogress")
        refreshContractGrid()
    },
    validatedClicked: function(){
        summaryVM.set("selected", "Validated")
        summaryVM.addhoverclass("classvalidated")
        refreshContractGrid()
    },
    resolutionClicked: function(){
        summaryVM.set("selected", "Resolution")
        summaryVM.addhoverclass("classresolution")
        refreshContractGrid()
    },
    refresh: function(){
        var self = this

        ajaxPost("/dashboard/getfilestatussummary", {}, function(res){
            self.set("new", res.Data.New)
            self.set("inProgress", res.Data.Progress)
            self.set("validated", res.Data.Validated)
            self.set("resolution", 0)
        })
    }
})

var contractGrid = null;
var renderContractGrid = function(){
    return $('#contract-grid').kendoGrid({
        dataSource: new kendo.data.DataSource({
            serverPaging: true,
            pageSize: 100,
            transport: {
                read: function(o){
                    ajaxPost("/dashboard/getuploadinformationtrue", {
                        statussummary: summaryVM.get("selected"),
                        page: o.data.page,
                        pageSize: o.data.pageSize,
                        skip: o.data.skip,
                        take: o.data.take
                    }, function(res){
                        o.success(res);
                    })
                }
            },
            schema: {
                parse: function(res){
                    return {
                        Data: res.Res.Data,
                        Total: res.Total
                    }
                },
                model: {
                    id: "Id",
                    fields: {
                        Filenamesource: { type: "string", editable: false },
                        // LibraryId: { type: "string" },
                        Cluster: { type: "string", editable: false },
                        ClusterConfident: { type: "number", editable: false },
                        SubCluster: { type: "string", editable: false },
                        SubClusterConfident: { type: "number", editable: false },
                        UploadTime: { type: "string", editable: false },
                        UploadBy: { type: "string", editable: false },
                    }
                },
                data: "Data",
                total: "Total"
            },
        }),
        columns: [
            { 
                command: [{ 
                    name: "destroy", 
                    text: " ", 
                    template: '<button data-command="destroy" class="btn btn-xs red default k-button-icontext k-grid-delete btn-delete-grid"><i class="fa fa-times"></i></button>',
                    width: 35,
                }],
                width: 25,
                headerAttributes: { style: "display: none" }
            },
            { field: 'ContractName', width: 100, title: 'Contract Name', headerAttributes: { colspan: 2 } },
            { field: 'ContractManager', width: 150, title: 'Contract Manager' },
            {
                title: 'Contract ID',
                field: 'ContractId',
                width: 70,
                headerAttributes: { style: "text-align: center" },
                attributes: { style: "text-align: left" },
            },
            {
                title: 'Machine Learning Classification',
                field: 'Cluster',
                width: 50,
                headerAttributes: { style: "text-align: center", colspan: 2 },
                attributes: {"class": "highlight" },
                template: function(e){
                    return '<span>' + e.ClusterDesc + '</span>'
                },
            },
            {
                title: "",
                headerAttributes: { style: "display: none" },
                attributes: {"class": "highlight"},
                width: 100,
                template: function(e){
                    return kendo.toString(e.ClusterConfident, "n0") + "% Confidence";
                }
            },
            {
                title: 'Machine Learning Type',
                field: 'SubCluster',
                width: 30,
                headerAttributes: { style: "text-align: center", colspan: 2 },
                attributes: {"class": "highlight" },
                template: function(e){
                    return '<span>' + e.SubClusterDesc + '</span>'
                },
            },
            {
                title: "",
                headerAttributes: { style: "display: none" },
                attributes: {"class": "highlight"},
                width: 100,
                template: function(e){
                    return kendo.toString(e.SubClusterConfident, "n0") + "% Confidence";
                }
            },
            {
                title: 'Upload Time',
                width: 100,
                template: function(e){
                    return moment(e.UploadTime).format("DD MMM YYYY HH:mm:ss")
                }
            }, 
            { 'field': 'UploadBy', width: 65, title: 'Uploaded By' },
            { 
                width: 120,
                template: function(e){ 
                    var label = function(){
                        switch(e.StatusSummary){
                            case "New Contract": return "Begin Review";
                            case "Progress": return "Continue Review";
                            case "Validated": return "View Review";
                        }
                    }
                    
                    return '<a href="/dashboard/review?id='+ e.Id + '" class="btn dark review">' + label() + '</a>'
                },
                attributes: { "class": "text-center" }
            },
        ],
        remove: function(e) {
            ajaxPost("/upload/deleteclassification", e.model.id, function(data){
                refreshContractGrid()
            })
        },
        dataBound: function(e){
            statussummaries = e.sender._data;
        },
        editable: true,
        height: 320,
        scrollable: {
            virtual: true
        },
    }).data('kendoGrid');
}

var refreshContractGrid = function(){
    contractGrid.dataSource.read();
    contractGrid.refresh();

    summaryVM.refresh()
}

var loadingModal = kendo.observable({
    progressValue: 0,
    message: "Uploading",
    isLoading: false,
    completed: false
})

kendo.bind($("#modal-loading"), loadingModal)

var tbdModal = kendo.observable({
    files: []
})
kendo.bind($("#modal-tbd"), loadingModal)

var uploadedId = ""

var createGuid = function() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1);
  }
  return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
    s4() + '-' + s4() + s4() + s4();
}

var gridTBDData = kendo.observable({
    data: [],
    deleted: []
});

var confirmAll = function() {
    var gridTbd = $("#grid-tbd").data("kendoGrid");

    _.each(gridTbd.dataSource.data(), function(v){ 
        v.set("isClusterConfirmed", true)
        v.set("isSubclusterConfirmed", true)
    })
}

var saveGridTBD = function(){
    var gridTbd = $("#grid-tbd").data("kendoGrid");
    var tbdDataSource = gridTbd.dataSource;
    var tbdDataSourceData = tbdDataSource.data()
    console.log("dataSource", tbdDataSource);

    // var libidTotal = 0;
    var contractNameTotal = 0
    var classificationTotal = 0
    var typeTotal = 0

    _.each(tbdDataSourceData, function(value){
        // if(value.get("LibraryId").length > 0)
        //     libidTotal++

        if(value.get("ContractName").length > 0)
            contractNameTotal++

        if(value.get("isClusterConfirmed"))
            classificationTotal++

        if(value.get("isSubclusterConfirmed"))
            typeTotal++
    });

    var totalData = tbdDataSource.total()
    if(
        // libidTotal == totalData && 
        contractNameTotal == totalData && classificationTotal == totalData && typeTotal == totalData){
        $('#modal-tbd').modal('hide');

        if (totalData > 0) {
            if ( ! loadingModal.get("isLoading")) {
                $('#modal-loading').modal('show');
                loadingModal.set("isLoading", true)
            }

            checkProgress("Done", function(){
                $('#modal-loading').modal('hide');
                loadingModal.set("isLoading", false)
                loadingModal.set("completed", false)

                setTimeout(function(){
                    loadingModal.set("message", "Uploading")
                }, 1000)
                
                refreshContractGrid()
                loadingModal.set("progressValue", 0)
            })
        }

        var findNonEditedData = function(v){
            return _.find(gridTBDData.get("data"), function(value){
                return value.Id == v.Id
            })
        }

        //rollback if cluster not confirmed (but edited)
        _.each(tbdDataSourceData, function(v){
            if( ! v.isClusterConfirmed)
                v.Cluster = findNonEditedData(v).Cluster;

            if( ! v.isSubclusterConfirmed)
                v.SubCluster = findNonEditedData(v).SubCluster;

            v.issaved = true;
        })

        ajaxPost("/upload/updateclassification", tbdDataSourceData, function(data){    
            refreshContractGrid()
        })

        _.each(gridTBDData.get("deleted"), function(value){
            ajaxPost("/upload/deleteclassification", value, function(data){
                refreshContractGrid()
            })
        });

        ajaxPost("/upload/extractinghtml", tbdDataSourceData, function(data){
            refreshContractGrid()
        })

        $('#modal-tbd').modal('hide');
        $('#grid-tbd').kendoGrid('destroy').empty();
    } else {
        alert("Please complete and confirm inputs")
    }
}

var checkProgress = function(endWhen, callback) {
    var toTitleCase = function (str){
        return str.replace(/\w\S*/g, function(txt){ return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase(); });
    }

    if( ! loadingModal.get("completed")){
        setTimeout(function(){
            ajaxPost("/dashboard/checkuploadstatus", { Id: uploadedId }, function(res){
                loadingModal.set("message", toTitleCase(res.Data[0].Status))
                loadingModal.set("progressValue", res.Data[0].Precentage)

                if (res.Data[0].Status == endWhen)
                    loadingModal.set("completed", true)

                if ( ! loadingModal.get("completed"))
                    checkProgress(endWhen, callback)
                else {
                    callback()
                    return
                }
            })
        }, 1000)
    } else {
        callback()
        return
    }
}

var renderGridTBD = function() {
    var gridTBD = $('#grid-tbd').kendoGrid({
        dataSource: new kendo.data.DataSource({
            serverPaging: true,
            pageSize: 100,
            transport: {
                read: function(o){
                    ajaxPost("/dashboard/getuploadinformation", {
                        uploadid: uploadedId,
                        page: o.data.page,
                        pageSize: o.data.pageSize,
                        skip: o.data.skip,
                        take: o.data.take
                    }, function(res){
                        console.log("-----------------------------", res);
                        o.success(res);
                        $('#modal-tbd').modal('show');
                    })
                }
            },
            schema: {
                parse: function(res){
                    var data = res.Res.Data
                    _.each(data, function(value){
                        value.ContractId = 0
                        value.isClusterConfirmed = false
                        value.isSubclusterConfirmed = false
                    });

                    gridTBDData.set("data", data)
                    console.log(">>>>>>>>>", data);
                    return {
                        Data: data,
                        Total: res.Total
                    }
                },
                model: {
                    id: "Id",
                    fields: {
                        ContractName: { type: 'string', editable: true },
                        ContractManager: { type: 'string', editable: false },
                        Cluster: { type: "string", editable: false },
                        ContractId: { type: "number", editable: true },
                        SubCluster: { type: "string", editable: false },
                        // LibraryId: { type: "string", editable: true },
                    }
                },
                data: "Data",
                total: "Total"
            },
        }),
        columns: [
            { 
                command: [{ 
                    name: "destroy", 
                    text: " ", 
                    template: '<button data-command="destroy" class="btn btn-xs red default k-button-icontext k-grid-delete"><i class="fa fa-times"></i></button>',
                    width: 35,
                }],
                width: 25, 
                headerAttributes: { style: "display: none" }
            },
            { field: 'ContractName', title: 'Contract Name', headerAttributes: { colspan: 2 } },
            { field: 'ContractManager', title: 'Contract Manager'  },
            // {
            //     title: 'Library Id',
            //     field: 'LibraryId',
            //     width: 120,
            //     headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align: center;" }
            // },
            {
                title: 'Scan Quality',
                width: 120,
                attributes: { class: 'text-center' },
                headerAttributes: {
                    style: "text-align: center",
                },
                attributes: {
                    style: "text-align:center;"
                },
                template: function(e){
                    return "Accept";
                }
            },
            {
                title: 'Contract ID',
                field: 'ContractId',
                width: 120,
                headerAttributes: {style: "text-align: center",},
                attributes: { style: "text-align: left" },
            },
            {
                title: 'Machine Learning Classification',
                field: 'Cluster',
                width: 120,
                headerAttributes: { style: "text-align: center", colspan: 4 },
                attributes: {"class": "highlight"},
                template:function(e){
                    var select;
                    var clusters = ['internal', 'external']
                    select = "<select id='id-cluster'>";
                    for(i  in clusters ){
                        var selected = '';
                        if(clusters[i] == e.Cluster){
                          var selected = 'selected'      
                        }
                        select += "<option "+selected+" value="+clusters[i]+">"+getClusterTitleCase(clusters[i])+"</option>"
                    }
                    select +="</select>";
                    return select
                }
            },
            {
                title: "",
                headerAttributes: { style: "display: none" },
                attributes: {"class": "highlight"},
                width: 100,
                template: function(e){
                    return kendo.toString(e.ClusterConfident, "n0") + "% Confidence";
                }
            },
            {
                title: " " ,
                // field: 'ClusterInfo',
                width: 35,
                headerAttributes: { style: "display: none" },
                attributes: {"class": "highlight"},
                template: function(e){
                    return '<span><i class="fa fa-info-circle information" title="' + kendo.toString(e.ClusterConfident, "n") + '% Confidence" style="font-size:15px; "></i></span>'
                },
            },
            { 
                command: [{ 
                    name: "cluster",
                    template: function(e){
                        return '<a role="button" class="k-button k-button-icontext k-grid-Confirm k-grid-cluster" href="#">'+'<span class=" "></span>Confirm'+'</a>'+
                        '<span class="k-plok-cluster">Confirmed <i class="fa fa-check-circle"></i></span>'
                    },
                    click: function(e) {
                        var dataItem = this.dataItem($(e.currentTarget).closest("tr"));
                        cluster = $('tr[data-uid="'+dataItem.uid+'"]').find('#id-cluster').val()
                        dataItem.Cluster = cluster.toLowerCase();
                        dataItem.set("isClusterConfirmed", true)
                    }
                }], 
                title: " " , width: 100, headerAttributes: { style: "display: none" },
                attributes: { "class": "highlight" },
            },
            {
                title: 'Machine Learning Type',
                field: 'SubCluster',
                width: 120,
                headerAttributes: { style: "text-align: center", colspan: 4 },
                attributes: {"class": "highlight"},
                template:function(e){
                    var select;
                    var subclusters = ['sla', 'msa', 'msla']
                    select = "<select id='id-subcluster'>";
                    for(i  in subclusters ){
                        var selected = '';
                        if(subclusters[i] == e.SubCluster){
                          var selected = 'selected'      
                        }
                        select += "<option "+selected+">"+subclusters[i].toUpperCase()+"</option>"
                    }
                    select +="</select>";
                    return select
                }
            },
            {
                title: "",
                headerAttributes: { style: "display: none" },
                attributes: {"class": "highlight"},
                width: 100,
                template: function(e){
                    return kendo.toString(e.SubClusterConfident, "n0") + "% Confidence";
                }
            },
            {
                title: " " ,
                // field: 'SubClusterInfo',
                width: 35,
                headerAttributes: { style: "display: none" },
                attributes: {"class": "highlight"},
                template: function(e){
                    return '<span><i class="fa fa-info-circle information" title="' + kendo.toString(e.SubClusterConfident, "n") + '% Confidence" style="font-size:15px;  "></i></span>'
                },
            },
            { 
                command: [{ 
                    name: "subcluster",
                    template: function(){
                        return '<a role="button" class="k-button k-button-icontext k-grid-Confirm k-grid-subcluster" href="#"><span class=" "></span>Confirm</a>'+
                        '<span class="k-plok-subcluster">Confirmed <i class="fa fa-check-circle"></i></span>'
                    },
                    click: function(e) {
                        var dataItem = this.dataItem($(e.currentTarget).closest("tr"));
                        subcluster = $('tr[data-uid="'+dataItem.uid+'"]').find('#id-subcluster').val()
                        dataItem.SubCluster =subcluster.toLowerCase();
                        dataItem.set("isSubclusterConfirmed", ( ! dataItem.get("isSubclusterConfirmed")))
                    } 
                }], 
                title: "" , width: 100, headerAttributes: { style: "display: none" },
                attributes: { "class": "highlight", css: "text-align: center;" },
            }
        ],
        remove: function(e) {
            console.log("Removing e", e);
            console.log("Removing e model", e.model);
            console.log("Removing e model name", e.model.name);
            gridTBDData.deleted.push(e.model.Id)
        },
        editable: true,
        scrollable: {
            virtual: true
        },
        dataBound: function(e) {
            $(".information").tooltipster({
                trigger: 'hover'
            });

            var grid = $("#grid-tbd").data("kendoGrid");
            var gridData = grid.dataSource.view();
            
            for (var i = 0; i < gridData.length; i++) {
                var currentUid = gridData[i].uid;

                var currenRow = grid.table.find("tr[data-uid='" + currentUid + "']");

                var clusterButton = $(currenRow).find(".k-grid-cluster");
                var clusterInfo = $(currenRow).find(".k-plok-cluster");

                if(gridData[i].isClusterConfirmed) {
                    clusterButton.hide();
                    clusterInfo.show();
                } else {
                    clusterButton.show();
                    clusterInfo.hide()
                }

                var subclusterButton = $(currenRow).find(".k-grid-subcluster");
                var subclusterInfo = $(currenRow).find(".k-plok-subcluster");
                
                if(gridData[i].isSubclusterConfirmed) {
                    subclusterButton.hide();
                    subclusterInfo.show();
                } else {
                    subclusterButton.show();
                    subclusterInfo.hide()
                }
            }

            // var sender = e.sender;
            // var cellToEdit = sender.content.find("td:eq(4), td:eq(7)");
            // sender.editCell(cellToEdit);
        }
    }).data('kendoGrid');
}

$(document).ready(function($) {
    kendo.bind($("#applicationStatus"), summaryVM)

    ajaxPost("/upload/deleteunsavedfile", {}, function(){})

    var urlGet = readGetFromUrl()
    if(urlGet != undefined)
        if(urlGet.p != undefined)
            var index = (urlGet.p[0] == "2") ? 'last' : 'first'

    contractGrid = renderContractGrid()

    var slider = $('.my-slider').unslider({
        index: index,
        nav: false,
        infinite: true,
        arrows: {
            prev: '<a class="unslider-arrow prev"><img src="/res/img/arrow-left.png" alt="prev" height="200px" width="30px"/></a>',
            next: '<a class="unslider-arrow next"><img src="/res/img/arrow-right.png" alt="next" height="200px" width="30px"/></a>',
        }
    });

    setTimeout(function(){
        slider.unslider('setIndex:1');
    }, 2000)

    $("#inputFile").kendoUpload({
        async: {
            saveUrl: "/upload/uploadhtmldata",
            removeUrl: "#",
            removeField: "fileNames[]",
            autoupload: false,
            batch: true
        },
        validation: {
            allowedExtensions: [".doc", ".docx", ".pdf", ".html"]
        },
        dropZone: ".dropZoneElement",
        success: function(e){
            uploadedId = e.response.Data
            if ( ! loadingModal.get("isLoading")) {
                $('#modal-loading').modal('show');
                loadingModal.set("isLoading", true)
            }

            checkProgress("Extracting Html", function(){
                $('#modal-loading').modal('hide');
                loadingModal.set("isLoading", false)
                loadingModal.set("completed", false)
                loadingModal.set("message", "Tagging Html")
                
                refreshContractGrid()
                renderGridTBD()
                loadingModal.set("progressValue", 0)
            })
        }
    });

    $(".k-button.k-upload-button").hide();
    $('.k-widget.k-upload.k-header').hide();

    $('.dropZoneElement').on('dragenter', function() {
        $('.dropZoneElement').has('div').find('p').css({'color': '#2b2d2f'})
        $('.dropZoneElement').css({'border': '2px dashed #2b2d2f'})
    });
    $('.dropZoneElement').on('dragleave', function() {
        $('.dropZoneElement').has('div').find('p').css({'color': '#c7c7c7'})
        $('.dropZoneElement').css({'border': '2px dashed #c7c7c7'})
    });
    $('.dropZoneElement .textWrapper').on('dragenter', function() {
        $('.dropZoneElement').has('div').find('p').css({'color': '#2b2d2f'})
        $('.dropZoneElement').css({'border': '2px dashed #2b2d2f'})
    });
    $('.dropZoneElement .textWrapper').on('dragleave', function() {
        $('.dropZoneElement').has('div').find('p').css({'color': '#c7c7c7'})
        $('.dropZoneElement').css({'border': '2px dashed #c7c7c7'})
    });

    // $("#contract-grid > .k-grid-content > .k-virtual-scrollable-wrap >table > tbody").on("mouseenter", "tr", function (e) {
    //     var thisuid = $(this).attr('data-uid')
    //     $.each(statussummaries, function(i,v){
    //         // statussummary.push({status:v.StatusSummary, uid:v.uid})
    //         if (thisuid == v.uid) {
    //             if (v.StatusSummary == "Progress") {
    //                 $('.classcontract').addClass('classnonhover');
    //                 $('.classprogress').addClass('classhover');
    //             }else if (v.StatusSummary == "Validated") {
    //                 $('.classcontract').addClass('classnonhover');
    //                 $('.classvalidated').addClass('classhover');
    //             }else if (v.StatusSummary == "Resolution") {
    //                 $('.classcontract').addClass('classnonhover');
    //                 $('.classresolution').addClass('classhover');
    //             };
    //         }
    //     })
    // });
    // $("#contract-grid > .k-grid-content > .k-virtual-scrollable-wrap >table > tbody").on("mouseleave", "tr", function (e) {
    //     $('.classcontract').removeClass('classnonhover').addClass('classhover');
    //     $('.classprogress').removeClass('classhover');
    //     $('.classvalidated').removeClass('classhover');
    //     $('.classresolution').removeClass('classhover');
    // });

    setWidth()
    summaryVM.refresh()
});

var setWidth = function(){
    $(".page-content").height($(window).height() - 105)
    $(".page-content").css("min-height", $(window).height() - 105);
    $(".dropZoneElement").height($(window).height() - 588)
    $(".dropZoneElement").width($(".my-slider").width())
}

$(window).resize(setWidth)