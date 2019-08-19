var typestock = {}

typestock.datatypestock = ko.observableArray([])
typestock.dataSearch = ko.observable("")
typestock.DatePageBar=ko.observable()
typestock.id = ko.observable()
model.ModeName = ko.observable("#Edittypestock");
typestock.button = ko.observable("Save")

typestock.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    typestock.DatePageBar(page)
}

typestock.getdatatypestock = function(callback) {
    model.Processing(true)
    ajaxPost('/master/getdatatypestock', {}, function(res) {
        if (res.Total === 0) {
            swal({title:"Error!",
            text: res.Message,
            type: "error",
            confirmButtonColor:"#3da09a"
        })
            return
        }
        model.Processing(false)
        typestock.datatypestock(res.Data)

        callback()
    }, function() {
        viewModel.isLoading(false)
        swal({
            title:"Error!", 
            text: "Unknown error, please try again", 
            type:"error"})
    })
}

typestock.renderGrid = function(search) {
    var data = typestock.datatypestock()
    if (search != undefined) {
        var results = _.filter(typestock.datatypestock(), function(item) {
            return item.name.indexOf(search) > -1 || item.code.indexOf(search) > -1
        });

        data = results
    }
    if (search == "") {
        data = typestock.datatypestock()
    }

    if (typeof $('#gridts').data('kendoGrid') !== 'undefined') {
        $('#gridts').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,

            pageSize: 25
        }))
        return
    }

    var columns = [
         {
            title: 'Action',
            width: 30,
            template: "<a href=\"javascript:typestock.Edit('#: _id #', '#:code#', '#:name#')\" data-target=\".Edittypestock\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" +
                "#if(1==1){#" +
                "<a href=\"javascript:typestock.Delete('#: _id #','#: code #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#"
            // template: "<a href=\"javascript:typestock.Edit('#: _id #', '#: code #', '#: name #')\" data-target=\"#AddNewModal\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" +
            //     "<a href=\"javascript:typestock.delete('#: _id #', '#: code #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>"
        },
        {
            field: 'code',
            title: 'Code',
            width: 100
        },
        {
            field: 'name',
            title: 'Name',
            width: 300
        }
    ]

    $('#gridts').kendoGrid({
        dataSource: data,
        sortable: true,
        height: 500,
        filterable: false,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "TypeStock", function(row, sheet){

            });
        },
    })
}
typestock.exportExcel = function () {
    $("#gridts").getKendoGrid().saveAsExcel();
}

typestock.SaveData = function() {
        var code = $("#code").val()
        var name = $("#name").val().trim().toUpperCase()
        if (name == "") {
            swal({
                title:"Error!",
                text: "All field must filled",
                type: "error",
                confirmButtonColor:"#3da09a"
            });
        }   
            var param = {
                ID : typestock.id(),
                code: code,
                name: name
            }
            var url = "/master/insertnewtypestock"
            model.Processing(true)
            ajaxPost(url, param, function(e) {
                code = $("#code").val("")
                name = $("#name").val("")
                model.Processing(false)
                if (e.status == true) {
                    swal({
                        title: "Success!", 
                        text: e.Message,
                        type: "success",
                        confirmButtonColor: '#3da09a',
                        confirmButtonText: "OK"});
                    typestock.init()
                    $('#AddNewModal').modal('hide');
                } else {
                    swal({
                        title:"Error!",
                        text: e.Message,
                        type: "error"});
                }

            })

       
}
typestock.Edit= function(id, code, name) {
    typestock.id(id)
    $("#code").val(code)
    $("#name").val(name)
    $("#AddNewModal").modal("show")
    typestock.button("Update")

}
typestock.Delete = function(id, code){
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + code + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "btn-danger",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function (res) {
        if (res) {
            var url = "/master/deletetypestock";
            var param = {
                ID: id
            };
            ajaxPost(url, param, function (data) {
                if (data.Message != "OK") {
                    swal({
                        title:'Warning',
                        text: data.Message,
                        type: 'error',
                        confirmButtonColor:"#3da09a"});
                    model.Processing(false);
                } else {
                    swal({
                        title:'Success',
                        text: 'Type Purchase has been deleted!',
                        type: 'success',
                        confirmButtonColor:"#3da09a"
                    });
                    typestock.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}
typestock.AddNew= function () {
    typestock.button("Save")
    typestock.id("")
    var code = $("#code").val("")
    var name = $("#name").val("")
}

typestock.CancelData = function() {
    var code = $("#code").val("")
    var name = $("#name").val("")
    if (code == "" || name == "") {
        swal({
            title:"Error!",
            text: "All field must filled",
            type: "error",
            confirmButtonColor:"#3da09a"});
        typestock.init()
    }
}

typestock.search = function() {
    var search = typestock.dataSearch()
    typestock.renderGrid(search)
},

typestock.init = function() {
    typestock.getdatatypestock(function() {
        typestock.renderGrid()
        typestock.getDateNow()
    })
}

typestock.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["code", "name"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridts").data("kendoGrid").dataSource.query({ filter: filter });
}

$(function() {
    typestock.init()
    $("#textSearch").on("keyup blur change", function () {
        typestock.filterText();
    });
})