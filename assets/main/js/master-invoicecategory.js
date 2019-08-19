var invcategory = {}

invcategory.newRecord = function() {
    return {
        _id: ko.observable(''),
        ID: ko.observable(''),
        Code: ko.observable(''),
        Name: ko.observable(''),
    }
}

invcategory.record = ko.mapping.fromJS(invcategory.newRecord())
invcategory.DatePageBar = ko.observable()
invcategory.datainvcategory = ko.observableArray([])
invcategory.id = ko.observable()
invcategory.text = ko.observable("Save")


invcategory.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    invcategory.DatePageBar(page)
}


invcategory.addNewModal = function() {
    invcategory.id("")
    $("#code").val("")
    $("#name").val("")
    invcategory.text("Save")
}

invcategory.editInvCategory = function(id, code, name) {
    invcategory.id(id)
    console.log(id, code, name)
    $("#code").val(code)
    $("#name").val(name)
    $("#AddNewModal").modal("show")
    invcategory.text("Update")
}

invcategory.SaveData = function() {
    var name = $("#name").val().trim()
    if (name == "") {
        swal({
            title:"Error!",
            text: "Category name is empty!",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    }else {
        var param = {
            ID: invcategory.id(),
            Name: name
        }
        console.log(param)
        var url = "/master/insertnewinvoicecategory"
        model.Processing(true)
        ajaxPost(url, param, function (e) {
            name = $("#name").val("")
            model.Processing(false)
            if (e.Status == true) {
                swal({
                        title: "Success!", 
                        text: e.Message,
                        type: "success",
                        confirmButtonColor: '#3da09a',
                        confirmButtonText: "OK"});
                invcategory.init()
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
}

invcategory.GetDataInvCategory = function(callback) {
    model.Processing(true)
    ajaxPost("/master/getdatainvoicecategory", {}, function(res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        model.Processing(false)
        invcategory.datainvcategory(res.Data)

        callback()
    }, function() {
        viewModel.isLoading(false)
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

invcategory.renderGrid = function(search) {
    var data = invcategory.datainvcategory()
    if (search != undefined) {
        var results = _.filter(invcategory.datainvcategory(), function(item) {
            return item.categoryname.indexOf(search) > -1 || item.categorycode.indexOf(search) > -1
        })
        data = results
        console.log(data)
    }
    if (search == "") {
        data = invcategory.datainvcategory()
    }
    if(typeof $("#gridinvcategory").data("kendoGrid") != 'undefined') {
        $("#gridinvcategory").data("kendoGrid").setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 15
        }))
        return
    }
    var columns = [
        {
            title: "Action",
            width: 35,
            template: "<a href=\"javascript:invcategory.editInvCategory('#: _id #', '#: categorycode #', '#: categoryname #')\" data-target=\".EditCustomer\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;"+ "#if(1==1){#" + "<a href=\"javascript:invcategory.deleteInvCategory('#: _id #', '#: categorycode #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#"
        }, {
            field: "categorycode",
            title: "Category Code",
            width: 100
        }, {
            field: "categoryname",
            title: "Category Name",
            width: 300
        }
    ]

    $('#gridinvcategory').kendoGrid({
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

invcategory.deleteInvCategory = function(id, code) {
    model.Processing(true)
    swal({
        title: "Are you sure to delete"+ code+"?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "btn-danger",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function(res) {
        if (res) {
            var url = "/master/deleteinvoicecategory";
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
    })
}

invcategory.exportExcel = function () {
    $("#gridinvcategory").getKendoGrid().saveAsExcel();
}

invcategory.search = function(txtsearch) {
    invcategory.renderGrid(txtsearch)
}

invcategory.init = function() {
    invcategory.GetDataInvCategory(function() {
        invcategory.renderGrid()
        invcategory.getDateNow()
    })
}

$(function() {
    invcategory.init()
})