var sales = {}

sales.DatePageBar = ko.observable()
sales.text = ko.observable("Save")
sales.filterindicator = ko.observable(false)
sales.dataMastersales = ko.observableArray([])
sales.dataMastersalesOriginal = ko.observableArray([])

sales.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    sales.DatePageBar(page)
}
sales.newRecord = function() {
    var page = {
        ID: "",
        SalesID: "",
        SalesName: "",
        Phone: "",
    }
    return page
}

sales.record = ko.mapping.fromJS(sales.newRecord())

sales.getDatasales = function(callback) {
    var param = {
        TextSearchCode: $("#filterCode").val(),
        TextSearchName: $("#filterName").val(),
    }

    ajaxPost('/master/getdatasales', param, function (res) {
        if (res.Total === 0) {
            swal({
                title:"Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        //console.log(res)
        sales.dataMastersales(res.Data)
        sales.dataMastersalesOriginal(res.Data)
        callback()
    }, function () {
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
    })
}


sales.reset = function() {
    ko.mapping.fromJS(sales.newRecord(), sales.record);
    if (sales.record._id)
        sales.record._id("");
    $(".formInput").val("")
    $("#qty").val(0)
}

sales.addNew = function() {
    sales.reset()
    $("#AddNewModal").modal("show")
    sales.text("Save")
}

sales.renderGrid = function() {
    var data = sales.dataMastersales();
    //console.log(data)
    if (typeof $('#gridsales').data('kendoGrid') !== 'undefined') {
        $('#gridsales').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
        }))
        return
    } 

    var columns = [{
        title: "Action",
        width: 90,
        template: "<button onclick='sales.edit(\"#: _id #\")' class='btn btn-xs btn-warning'><i class='fa fa-pencil'></i></button> <button onclick='sales.delete(\"#: _id #\",\"#: SalesID #\")' class='btn btn-xs btn-danger'><i class='glyphicon glyphicon-trash' aria-hidden='true'></i></button>"
    },{
        field: "SalesID",
        title: "Sales Code",
        width: 160,
    },{ 
        field: "SalesName",
        title: "Sales Name",
        width: 160,
    },{
        field: "Phone",
        title: "Phone",
        width: 60,
    }]
    $('#gridsales').kendoGrid({
        dataSource: {
            data: data,
        },
        height: 500,
        width: 140,
        sortable: true,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "Sales", function(row, sheet){
                
            });
        },
    })
}
sales.exportExcel = function () {
    $("#gridsales").getKendoGrid().saveAsExcel();
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

sales.saveData = function() {
    var change = ko.mapping.toJS(sales.record)   
    change.Qty = parseInt(change.Qty)
    var param = {
        Data: change,
    }

    change.SalesName = change.SalesName.toUpperCase().trim()
    var name = change.SalesName
    var vldName = FilterAlphanum(name)
    if (vldName == "") {
        return swal({
            title:'Warning!', 
            text:"You haven't fill the Sales Name", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }
    if (vldName.length < 3) {
        return swal({
            title:'Warning!', 
            text:"The Sales Name is too short (min 3 letters)", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }
    change.Phone = "" + change.Phone;

    var url = "/master/insertnewsales"
    swal({
        title: "Are you sure?",
        text: "You will submit this sales",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function(isConfirm) {
        if (isConfirm) {
           model.Processing(true)
            ajaxPost(url, param, function(e) {
                //console.log(e.IsError)
                if (e.IsError) {
                    setTimeout(function () {
                        swal({
                            title: "Warning!",
                            text: e.Message,
                            type: "error",
                            confirmButtonColor:"#3da09a"
                        }, function () {
                            model.Processing(false);
                        });
                    }, 200)
                } else {
                    swal({
                    title:'Success', 
                    text:'sales has been saved!',
                    type:'success',
                    confirmButtonColor:"#3da09a"
                });
                    sales.init()
                    sales.reset()
                    $('#AddNewModal').modal('hide');
                    model.Processing(false);
                }
                model.Processing(false)
            })
        } else {
            swal({
                title:"Cancelled",
                type:"error",
                confirmButtonColor:"#3da09a"
            });
        }
    });
}

sales.edit = function(id) {    
    var data = _.find(sales.dataMastersales(), function(o) {
        return o._id == id;
    });
    ko.mapping.fromJS(data,sales.record)
    $("#AddNewModal").modal("show")
    sales.text("Update")
}

sales.delete = function(id, INVID) {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + INVID + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "#3da09a",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function (res) {
        if (res) {
            var url = "/master/deletesales";
            var param = {
                ID: id
            };
            ajaxPost(url, param, function (data) {
                if (data.Message != "OK") {
                    swal({
                        title:'Warning',
                        text: data.Message,
                        type: 'error',
                        confirmButtonColor:"#3da09a"
                    });
                    model.Processing(false);
                } else {
                    swal({
                        title:'Success',
                        text: 'sales has been deleted!',
                        type: 'success',
                        confirmButtonColor:"#3da09a"
                    });
                    sales.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}


sales.init = function() {
    sales.getDatasales(function(){
        sales.renderGrid()  
        sales.getDateNow()
    })
}

$(function(){
    sales.init()
})