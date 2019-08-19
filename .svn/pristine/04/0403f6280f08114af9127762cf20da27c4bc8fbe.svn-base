var gridVM = kendo.observable({
    dataSource: new kendo.data.DataSource({
        serverPaging: true,
        pageSize: 100,
    	transport: {
            read: function(o){
                ajaxPost("/refvendors/getentities", {
                    page: o.data.page,
                    pageSize: o.data.pageSize,
                    skip: o.data.skip,
                    take: o.data.take,
                }, function(res){
                    o.success(res);
                    console.log(res)
                })
            },
            create: function(o){
                var griddata = $("#vendor-grid").getKendoGrid().dataSource.data()
                var dirty = $.grep(griddata, function(item) {
                    return item.dirty
                });

                var data = {}
                var datas = []
                if (dirty.length != 0){
                    for (var i = 0 ; i < dirty.length ; i++){
                        data.Code2 = dirty[i].Code2
                        data.Code3 = dirty[i].Code3
                        data.Name = dirty[i].Name
                        datas.push(data)
                    }   
                }

                var param = datas
                console.log(datas)

                ajaxPost("/refvendors/saveentities", param, function(data){
                    swal("Success!", data.Message, "success");
                })
            },
            update: function(o){
                var griddata = $("#vendor-grid").getKendoGrid().dataSource.data()
                var dirty = $.grep(griddata, function(item) {
                    return item.dirty
                });

                var data = {}
                var datas = []
                if(dirty.length != 0){
                    for (var i = 0 ; i < dirty.length ; i++){
                        data.Id = dirty[i].Id
                        data.Code2 = dirty[i].Code2
                        data.Code3 = dirty[i].Code3
                        data.Name = dirty[i].Name
                        datas.push(data)
                    }   
                }

                var param = datas
                console.log(datas)

                ajaxPost("/refvendors/updateentities",param,function(data){
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
                    Code2: { type: "string" },
                    Code3: { type: "string" },
                    Name: { type: "string" },
                }
            },
            data: "Data",
            total: "Total"
        },
    })
});

var grid = $('#vendor-grid').kendoGrid({
    dataSource: gridVM.dataSource,
    columns: [
        { field: 'Code2', title: 'Code 2' },
        { field: 'Code3', title: 'Code 3' },
        { field: 'Name', title: 'Name' },
    ],
    height: 300,
    filterable: true,
    toolbar: ["create", "save", "cancel"],
    editable: true,
    navigatable:true,
    scrollable: {
        virtual: true
    },
}).data('kendoGrid');