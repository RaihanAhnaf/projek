var clauseGrid = kendo.observable({
    source: [ 
        { val: "internal-sla", txt: "Internal SLA" },
        { val: "internal-msa", txt: "Internal MSA" },
        { val: "external-sla", txt: "External SLA" },
        { val: "external-msa", txt: "External MSA" },
    ],
    filter: "",
})
clauseGrid.changed = function(){
    var filter = clauseGrid.get("filter")

    var splittedFilter = clauseGrid.get("filter").split("-")
    clauseGridVM.set("cluster", splittedFilter[0])
    clauseGridVM.set("subcluster", splittedFilter[0] + splittedFilter[1]) 

    $('#clause-grid').data('kendoGrid').dataSource.read();
    $('#clause-grid').data('kendoGrid').refresh();
}

var clauseGridVM = kendo.observable({
    cluster: "",
    subcluster: "",
    dataSource: new kendo.data.DataSource({
        serverPaging: true,
        pageSize: 100,
    	transport: {
            read: function(o){
                ajaxPost("/refdocumentclause/getclause", {
                    cluster: clauseGridVM.get("cluster"),
                    subcluster: clauseGridVM.get("subcluster"),
                    page: o.data.page,
                    pageSize: o.data.pageSize,
                    skip: o.data.skip,
                    take: o.data.take,
                }, function(res){
                    o.success(res);
                })
            },
            create: function(o){
                var griddata = $("#clause-grid").getKendoGrid().dataSource.data()
                var dirty = $.grep(griddata, function(item) {
                    return item.dirty
                });

                var data = {}
                var datas = []
                if(dirty.length != 0){
                    for (var i = 0 ; i < dirty.length ; i++){
                        data.Cluster = dirty[i].Cluster
                        data.SubCluster = dirty[i].SubCluster
                        data.ClauseId = dirty[i].ClauseId
                        data.ClauseLabel = dirty[i].ClauseLabel
                        data.Required = dirty[i].Required
                        datas.push(data)
                    }   
                }

                var param = datas
                console.log(datas)

                ajaxPost("/refdocumentclause/savedata",param,function(data){
                    swal("Success!", data.Message, "success");
                })


            },
            update: function(o){
                var griddata = $("#clause-grid").getKendoGrid().dataSource.data()
                var dirty = $.grep(griddata, function(item) {
                    return item.dirty
                });

                var data = {}
                var datas = []
                if(dirty.length != 0){
                    for (var i = 0 ; i < dirty.length ; i++){
                        data.Id = dirty[i].Id
                        data.Cluster = dirty[i].Cluster
                        data.SubCluster = dirty[i].SubCluster
                        data.ClauseId = dirty[i].ClauseId
                        data.ClauseLabel = dirty[i].ClauseLabel
                        data.Required = dirty[i].Required
                        datas.push(data)
                    }   
                }

                var param = datas
                console.log(datas)

                ajaxPost("/refdocumentclause/updatedata",param,function(data){
                    swal("Success!", data.Message, "success");
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
                    Cluster: { type: "string" },
                    SubCluster: { type: "string" },
                    ClauseId: { type: "string" },
                    ClauseLabel: { type: "string" },
                    Required: { type: "string" },
                }
            },
            data: "Data",
            total: "Total"
        },
    })
});

var selectedCluster = kendo.observable({
    val: ""
})

var grid = $('#clause-grid').kendoGrid({
    dataSource: clauseGridVM.dataSource,
    columns: [
        { field: 'ClauseId', title: 'Clause Id' },
        { field: 'ClauseLabel', title: 'Clause Label' },
        { field: 'Required', title: 'Required' },
        
    ],
    height: 500,
    filterable: true,
    toolbar: ["create", "save", "cancel"],
    editable: true,
    scrollable: {
        virtual: true
    },
}).data('kendoGrid');

$("#filter").kendoDropDownList({
    dataValueField: 'val', 
    dataTextField: 'txt',
    optionLabel: "Select cluster..."
})

kendo.bind($("#clauseGrid"), clauseGrid)