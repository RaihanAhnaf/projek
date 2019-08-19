var typepurchase = {}

typepurchase.datatypepurchase = ko.observableArray([])
typepurchase.dataSearch = ko.observable("")
typepurchase.DatePageBar=ko.observable()
typepurchase.id = ko.observable()
model.ModeName = ko.observable("#EditTypePurchase");
typepurchase.button = ko.observable("Save")

typepurchase.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    typepurchase.DatePageBar(page)
}

typepurchase.getdataTypepurchase = function(callback) {
    model.Processing(true)
    ajaxPost('/master/getdatatypepurchase', {}, function(res) {
        if (res.Total === 0) {
            swal({title:"Error!",
            text: res.Message,
            type: "error",
            confirmButtonColor:"#3da09a"
        })
            return
        }
        model.Processing(false)
        typepurchase.datatypepurchase(res.Data)

        callback()
    }, function() {
        viewModel.isLoading(false)
        swal({
            title:"Error!", 
            text: "Unknown error, please try again", 
            type:"error"})
    })
}

typepurchase.renderGrid = function(search) {
    var data = typepurchase.datatypepurchase()
    if (search != undefined) {
        var results = _.filter(typepurchase.datatypepurchase(), function(item) {
            return item.name.indexOf(search) > -1 || item.code.indexOf(search) > -1
        });

        data = results
    }
    if (search == "") {
        data = typepurchase.datatypepurchase()
    }

    if (typeof $('#gridtp').data('kendoGrid') !== 'undefined') {
        $('#gridtp').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,

            pageSize: 25
        }))
        return
    }

    var columns = [
         {
            title: 'Action',
            width: 50,
            template: "<a href=\"javascript:typepurchase.Edit('#: _id #', '#:code#', '#:name#')\" data-target=\".EditTypePurchase\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" +
                "#if(1==1){#" +
                "<a href=\"javascript:typepurchase.Delete('#: _id #','#: code #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#"
            // template: "<a href=\"javascript:typepurchase.Edit('#: _id #', '#: code #', '#: name #')\" data-target=\"#AddNewModal\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" +
            //     "<a href=\"javascript:typepurchase.delete('#: _id #', '#: code #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>"
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

    $('#gridtp').kendoGrid({
        dataSource: data,
        sortable: true,
        height: 500,
        filterable: false,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "TypePurchase", function(row, sheet){

            });
        },
    })
}
typepurchase.exportExcel = function () {
    $("#gridtp").getKendoGrid().saveAsExcel();
}
function FilterAlphanum(str, charset) {
    if (!str) return str;
    charset = charset || "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    var re = "";
    for(var i=0; i < str.length; i++) {
        if (charset.includes(str.substr(i,1))) {    
            re += str.substr(i,1);
        }
    }
    return re;
}


typepurchase.SaveData = function() {
        var code = $("#code").val()
        var name = $("#name").val().trim().toUpperCase()
        var vldName = FilterAlphanum(name)
        if (vldName == "") {
            return swal({
                title:'Warning!', 
                text:"You haven't fill the Type Purchase Name", 
                type:"info",
                confirmButtonColor:"#3da09a"
            })
        }
        if (vldName.length < 3) {
            return swal({
                title:'Warning!', 
                text:"The Type Purchase Name is too short (min 3 letters)", 
                type:"info",
                confirmButtonColor:"#3da09a"
            })
        }
        if (name == "") {
            swal({
                title:"Error!",
                text: "All field must filled",
                type: "error",
                confirmButtonColor:"#3da09a"
            });
        }  
            var param = {
                ID : typepurchase.id(),
                code: code,
                name: name
            }
            var url = "/master/insertnewtypepurchase"
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
                    typepurchase.init()
                    $('#AddNewModal').modal('hide');
                } else {
                    swal({
                        title:"Error!",
                        text: e.Message,
                        type: "error"});
                }

            })

}
typepurchase.Edit= function(id, code, name) {
    typepurchase.id(id)
    $("#code").val(code)
    $("#name").val(name)
    $("#AddNewModal").modal("show")
    typepurchase.button("Update")

}
typepurchase.Delete = function(id, code){
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
            var url = "/master/deletetypepurchase";
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
                    typepurchase.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}
typepurchase.AddNew= function () {
    typepurchase.button("Save")
    typepurchase.id("")
    var code = $("#code").val("")
    var name = $("#name").val("")
}

typepurchase.CancelData = function() {
    var code = $("#code").val("")
    var name = $("#name").val("")
    if (code == "" || name == "") {
        swal({
            title:"Error!",
            text: "All field must filled",
            type: "error",
            confirmButtonColor:"#3da09a"});
        typepurchase.init()
    }
}

typepurchase.search = function() {
    var search = typepurchase.dataSearch()
    typepurchase.renderGrid(search)
},

typepurchase.init = function() {
    typepurchase.getdataTypepurchase(function() {
        typepurchase.renderGrid()
        typepurchase.getDateNow()
    })
}

typepurchase.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["code", "name"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridtp").data("kendoGrid").dataSource.query({ filter: filter });
}

$(function() {
    typepurchase.init()
    $("#textSearch").on("keyup blur change", function () {
        typepurchase.filterText();
    });
})