var category = {}

category.newRecord = function () {
    return {
        _id: ko.observable(''),
        ID: ko.observable(''),
        Code: ko.observable(''),
        Name: ko.observable(''),
    }
}
category.record = ko.mapping.fromJS(category.newRecord())
category.dataCategory = ko.observableArray([])
category.DatePageBar = ko.observable()
category.text = ko.observable("Save")
category.id = ko.observable()

category.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    category.DatePageBar(page)
}

category.getDataCategory = function (callback) {
    model.Processing(true)
    ajaxPost('/master/getdatacategory', {}, function (res) {
        if (res.Total === 0) {
            swal({
                title:"Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
        })
            return
        }
        model.Processing(false)
        category.dataCategory(res.Data)

        callback()
    }, function () {
        viewModel.isLoading(false)
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

category.renderGrid = function (search) {
    var data = category.dataCategory()
    if (search != undefined) {
        var results = _.filter(category.dataCategory(), function (item) {
            return item.name.indexOf(search) > -1 || item.code.indexOf(search) > -1
        });

        data = results
    }
    if (search == "") {
        data = category.dataCategory()
    }

    if (typeof $('#gridcategory').data('kendoGrid') !== 'undefined') {
        $('#gridcategory').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,

            pageSize: 25
        }))
        return
    }

    var columns = [
     {
        title: 'Action',
        width: 35,
        template: "<a href=\"javascript:category.editCategory('#: _id #', '#: code #', '#: name #')\" data-target=\".EditCustomer\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" + "#if(1==1){#" + "<a href=\"javascript:category.deleteCategory('#: _id #', '#: code #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#"
    },{
        field: 'code',
        title: 'Asset Code',
        width: 100
    }, {
        field: 'name',
        title: 'Asset Name',
        width: 300
    }]

    $('#gridcategory').kendoGrid({
        dataSource: data,
        pageable: false,
        sortable: true,
        height: 500,
        filterable: false,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "Category", function(row, sheet){

            });
        },
    })
}

category.exportExcel = function () {
    $("#gridcategory").getKendoGrid().saveAsExcel();
}

category.addNewModal = function() {
    category.id("")
    $("#code").val("")
    $("#name").val("")
    category.text("Save")
}

category.editCategory = function(id, code, name) {
    category.id(id)
    $("#code").val(code)
    $("#name").val(name)
    $("#AddNewModal").modal("show")
    category.text("Update")
}

category.deleteCategory = function(id, code) {
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
            var url = "/master/deletecategory";
            var param = {
                ID: id
            };
            ajaxPost(url, param, function (data) {
                if (data != "") {
                    swal({
                        title:'Warning',
                        text: data, 
                        type:'error',
                        confirmButtonColor:"#3da09a"
                    });
                    model.Processing(false);
                } else {
                    swal({
                        title:'Success',
                        text: 'User has been deleted!',
                        type: 'success',
                        confirmButtonColor:"#3da09a"
                    });
                    category.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}


category.SaveData = function () {
    var name = $("#name").val().trim()
    if (name == "") {
        swal({
            title:"Error!",
            text: "Asset Name is empty!",
            type: "error",
            confirmButtonColor:"#3da09a"
        });
    } else {
        var param = {
            ID: category.id(), 
            name: name
        }
        var url = "/master/insertnewcategory"
        model.Processing(true)
        ajaxPost(url, param, function (e) {
            name = $("#name").val("")
            model.Processing(false)
            if (e.status == true) {
                swal({
                        title: "Success!", 
                        text: e.Message,
                        type: "success",
                        confirmButtonColor: '#3da09a',
                        confirmButtonText: "OK"});
                category.init()
                $('#AddNewModal').modal('hide');
            } else {
                swal({
                    title:"Error!",
                    text: e.Message, 
                    type:"error",
                    confirmButtonColor:"#3da09a"
                });
            }
        })

    }
},

category.CancelData = function () {
    var code = $("#code").val("")
    var name = $("#name").val("")
    if (code == "" || name == "") {
        swal({
            title:"Error!",
            text: "All field must filled", 
            type:"error",
            confirmButtonColor:"#3da09a"
        });
        category.init()
    }
}

category.search = function (textsearch) {
    category.renderGrid(textsearch)
},

category.init = function () {
    category.getDataCategory(function () {
        category.renderGrid()
        category.getDateNow()
    })
}

$(function () {
    category.init()
})