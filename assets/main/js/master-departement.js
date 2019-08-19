var departement = {}

departement.newRecord = function () {
    return {
        _id: ko.observable(''),
        ID: ko.observable(''),
        Code: ko.observable(''),
        Name: ko.observable(''),
    }
}
departement.record = ko.mapping.fromJS(departement.newRecord())
departement.datadepartement = ko.observableArray([])
departement.DatePageBar = ko.observable()
departement.text = ko.observable("Save")
departement.id = ko.observable()

departement.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    departement.DatePageBar(page)
}

departement.getDatadepartement = function (callback) {
    model.Processing(true)
    ajaxPost('/transaction/getdatadepartment', {}, function (res) {
        if (res.IsError === true) {
            swal({
                title:"Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
        })
            return
        }
        model.Processing(false)
        departement.datadepartement(res.Data)
        //console.log(res)

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

departement.renderGrid = function (search) {
    var data = departement.datadepartement()
    // if (search != undefined) {
    //     var results = _.filter(departement.datadepartement(), function (item) {
    //         return item.name.indexOf(search) > -1 || item.code.indexOf(search) > -1
    //     });

    //     data = results
    // }
    // if (search == "") {
    //     data = departement.datadepartement()
    // }

    if (typeof $('#griddepartement').data('kendoGrid') !== 'undefined') {
        $('#griddepartement').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,

            pageSize: 25
        }))
        return
    }

    var columns = [
     {
        field: 'DepartmentCode',
        title: 'Departement Code',
        width: 300
    },{
        field: 'DepartmentName',
        title: 'Departement Name',
        width: 300
    }]

    $('#griddepartement').kendoGrid({
        dataSource: data,
        pageable: false,
        sortable: true,
        height: 500,
        filterable: false,
        scrollable: true,
        columns: columns,
        
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "Department", function(row, sheet){
                
            });
        },
    })
}


departement.exportExcel = function () {
    $("#griddepartement").getKendoGrid().saveAsExcel();
}

departement.addNewModal = function() {
    departement.id("")
    $("#name").val("")
    departement.text("Save")
}

departement.editdepartement = function(id, name) {
    departement.id(id)
    $("#name").val(name)
    $("#AddNewModal").modal("show")
    departement.text("Update")
}

departement.deletedepartement = function(id, name) {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + name + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "btn-danger",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function (res) {
        if (res) {
            var url = "/transaction/deletedepartement";
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
                    departement.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
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

departement.SaveData = function () {
    var name = $("#name").val()
    var vldName = FilterAlphanum(name)
    if (vldName == "") {
        return swal({
            title:'Warning!', 
            text:"You haven't fill the Department Name", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }else if (vldName.length < 3) {
        return swal({
            title:'Warning!', 
            text:"The Department Name is too short (min 3 letters)", 
            type:"info", 
            confirmButtonColor:"#3da09a"
        })
    } else {
        model.Processing(true)
        var url = "/transaction/savenewdepartment"
        var param ={
            DepartmentName : name.toUpperCase()
        }
        ajaxPost(url, param, function (data) {
            $("#AddNewModal").modal("hide")
            model.Processing(false)
            swal({
                title:"Success!",
                text: "Success Add New Data",
                type: "success",
                confirmButtonColor: "#3da09a"
            })
            departement.getDatadepartement(function () {
                departement.renderGrid()
                departement.getDateNow()
            })
        })
    }
},

departement.CancelData = function () {
    var name = $("#name").val("") 
    if (name == "") {
        swal({
            title:"Error!",
            text: "All field must filled", 
            type:"error",
            confirmButtonColor:"#3da09a"
        });
        departement.init()
    }
}

departement.search = function (textsearch) {
    departement.renderGrid(textsearch)
},

departement.init = function () {
    departement.getDatadepartement(function () {
        departement.renderGrid()
        departement.getDateNow()
    })
}

$(function () {
    departement.init()
})