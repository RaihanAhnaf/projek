var subclusterGridVM = kendo.observable({
    dataSource: new kendo.data.DataSource({
        serverPaging: true,
        pageSize: 100,
    	transport: {
            read: function(o){
                ajaxPost("/refmaster/getdatamaster", {
                    page: o.data.page,
                    pageSize: o.data.pageSize,
                    skip: o.data.skip,
                    take: o.data.take,
                    category: "DOC-SUBCLUSTER"
                }, function(res){
                    o.success(res);
                    console.log(res)
                })
            },
            create: function(o){
                var griddata = $("#subcluster-grid").getKendoGrid().dataSource.data()
                var dirty = $.grep(griddata, function(item) {
                    return item.dirty
                });

                var data = {}
                var datas = []
                if(dirty.length != 0){
                    for (var i = 0 ; i < dirty.length ; i++){
                        data.Category = "DOC-SUBCLUSTER"
                        data.Code = dirty[i].Code
                        data.Name = dirty[i].Name
                        data.Description = dirty[i].Description
                        data.Parent = dirty[i].ParentId
                        datas.push(data)
                    }   
                }

                var param = datas

                ajaxPost("/refmaster/savedata",param,function(data){
                    swal("Success!", data.Message, "success");
                })

            },
            update: function(o){
                var griddata = $("#subcluster-grid").getKendoGrid().dataSource.data()
                var dirty = $.grep(griddata, function(item) {
                    return item.dirty
                });

                var data = {}
                var datas = []
                if(dirty.length != 0){
                    for (var i = 0 ; i < dirty.length ; i++){
                        data.Id = dirty[i].Id
                        data.Category = "DOC-SUBCLUSTER"
                        data.Code = dirty[i].Code
                        data.Name = dirty[i].Name
                        data.Description = dirty[i].Description
                        data.Parent = dirty[i].ParentId
                        datas.push(data)
                    }   
                }

                var param = datas

                ajaxPost("/refmaster/updatedata",param,function(data){
                    swal("Success!", data.Message, "success");
                })
            },
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
                    Category: { type: "string" },
                    Code: { type: "string" },
                    Name: { type: "string" },
                    Description: { type: "string" },
                    Parent: { type: "number" },
                }
            },
            data: "Data",
            total: "Total"
        },
    })
});

var grid = $('#subcluster-grid').kendoGrid({
    dataSource: subclusterGridVM.dataSource,
    columns: [
     	{ field: 'Category', title: 'Category' },
        { field: 'Code', title: 'Code' },
        { field: 'Name', title: 'Name' },
        { field: 'Description', title: 'Description' },
        { 
            field: 'ParentId', 
            title: 'Parent',
            editor: function (container, options) {
                $('<input required name="' + options.field + '"/>')
                .appendTo(container)
                    .kendoDropDownList({
                        autoBind: false,
                        dataSource: {
                            transport:{
                                read : function(o){
                                    ajaxPost("/refmaster/getparentid", {}, function(res){
                                        o.success(res);
                                    })
                                }
                            }
                        }
                    });
                }
         },
    ],
    height: 300,
    filterable: true,
    toolbar: ["create", "save", "cancel"],
    navigatable: true,
    editable: true,
    scrollable: {
        virtual: true
    },

}).data('kendoGrid');


