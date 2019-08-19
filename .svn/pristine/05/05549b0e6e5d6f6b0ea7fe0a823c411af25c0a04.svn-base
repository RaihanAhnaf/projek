var inventory = {}

inventory.DatePageBar = ko.observable()
inventory.text = ko.observable("Save")
inventory.filterindicator = ko.observable(false)
inventory.dataMasterInventory = ko.observableArray([])
inventory.dataMasterInventoryOriginal = ko.observableArray([])
inventory.DataTypeStock = ko.observableArray([])
inventory.valueUnit = ko.observable("")
inventory.valueLocationSearch = ko.observable(0)
inventory.dataDropDownLocation = ko.observableArray([])
inventory.textSearch = ko.observable("")

inventory.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    inventory.DatePageBar(page)
}

inventory.newRecord = function() {
    var page = {
        ID: "",
        INVID: "",
        INVDesc: "",
        Unit: "",
        Type:""
    }
    return page
}

inventory.record = ko.mapping.fromJS(inventory.newRecord())

inventory.getDataLocation = function(callback) {
    $.ajax({
        url: "/transferorder/getuserlocations",
        success: function (json) {           
            inventory.dataDropDownLocation.removeAll();
            inventory.dataDropDownLocation(json);
            if (json.length > 0)
            {
                inventory.valueLocationSearch(json[0].LocationID);
                inventory.valueLocationSearch.valueHasMutated();
            }
            if (typeof callback == "function") callback();
        }
    });
}

inventory.getDefaultLocationID = function() {
    var locData = inventory.dataDropDownLocation()
    if (locData.length > 0) return locData[0].LocationID;
    return -1;
}

inventory.getDataInventory = function(callback) {    
    ajaxPost('/master/getalldatainventory', {
        LocationID: parseInt(inventory.valueLocationSearch())
        // LocationID: parseInt(inventory.valueLocationSearch()) || inventory.getDefaultLocationID()
    }, function(res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        // console.log(res)
        inventory.dataMasterInventory(res.Data)
        inventory.dataMasterInventoryOriginal(res.Data)
        callback()
    }, function() {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

inventory.reset = function() {
    $(".formInput").val("")
    $("#qty").val(0)
}

inventory.addNew = function() {
    inventory.reset()
    inventory.text("Save")
}

// inventory.search = function() {
//     inventory.getDataInventory(function() {
//         inventory.renderGrid();
//     });
// }

inventory.renderGrid = function(search) {
    var data = inventory.dataMasterInventory();
    if (search != undefined) {
        var results = _.filter(data, function (item) {
            var res = search.toUpperCase();
            return item.INVDesc.indexOf(res) > -1 || item.INVID.indexOf(res) > -1 || item.StoreLocationName.indexOf(res) > -1
        });
        data = results
    }
    if (search == "") {
        data
    }
    var columns = [
        // {
        //     title: "Action",
        //     width: 90,
        //     template: "<button onclick='inventory.edit(\"#: _id #\")' class='btn btn-xs btn-warning'><i class='fa fa-pencil'></i></button> <button onclick='inventory.delete(\"#: _id #\",\"#: INVID #\")' class='btn btn-xs btn-danger'><i class='glyphicon glyphicon-trash' aria-hidden='true'></i></button>"
        // }, 
        {
            field: "INVID",
            title: "Code",
            width: 100,
        }, {
            field: "INVDesc",
            title: "Name",
            width: 100,
        }, {
            field: "StoreLocationName",
            title: "Storehouse",
            width: 60,
        }, {
            field: "Unit",
            title: "Unit",
            width: 60,
        }, {
            field: "Type",
            title: "Type",
            width: 60,
            template: function(e){
                return e.Type;
            },
        },{
            field: "Beginning",
            title: "Beginning",
            width: 60,
        }, {
            field: "InInventory",
            title: "PO",
            width: 60,
        }, {
            field: "CMVInventory",
            title: "CMV",
            width: 60,
        }, {
            field: "OutInventory",
            title: "INV",
            width: 60,
        }, {
            field: "CMInventory",
            title: "CM",
            width: 60,
        }, {
            field: "TSInventory",
            title: "TS",
            width: 60,
        }, {
            field: "TRInventory",
            title: "TR",
            width: 60,
        }, {
            field: "Saldo",
            title: "Saldo",
            width: 60,
        }, {
            field: "UnitCost",
            title: "UnitCost",
            width: 60,
            attributes: {
                style: "text-align:right;"
            },
            template: "#=ChangeToRupiah(UnitCost)#",
        }, {
            field: "Total",
            title: "Total",
            width: 60,
            attributes: {
                style: "text-align:right;"
            },
            template: "#=ChangeToRupiah(Total)#",
        }, {
            field: "LastDate",
            title: "LastDate",
            width: 60,
            template:"#=moment(LastDate).format('DD-MMMM-YYYY')#"
        }

    ]
    $('#gridinventory').kendoGrid({
        dataSource: {
            data: data,
            pageSize: 10,
        },
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        height: 500,
        width: 140,
        sortable: true,
        scrollable: true,
        columns: columns,
    })
}

inventory.saveData = function() {
    var change = ko.mapping.toJS(inventory.record)
    change.Unit = inventory.valueUnit()
    change.INVDesc = change.INVDesc.trim()
    if (change.INVDesc == "") {
        return swal({
            title:'Warning!', 
            text:"You haven't fill the Description", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    };
    var param = {
        Data: change,
    }
    var url = "/master/saveinventory"
    swal({
        title: "Are you sure?",
        text: "You will submit this Inventory",
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
                    setTimeout(function() {
                        swal({
                            title: "Warning!",
                            text: e.Message,
                            type: "error",
                            confirmButtonColor: "#3da09a"
                        }, function() {
                            model.Processing(false);
                        });
                    }, 200)
                } else {
                    swal({
                        title: 'Success',
                        text: 'Inventory has been saved!',
                        type: 'success',
                        confirmButtonColor: "#3da09a"
                    });
                    inventory.init()
                    $('#AddNewModal').modal('hide');
                    model.Processing(false);
                }
                model.Processing(false)
            })
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}

inventory.edit = function(id) {
    var data = _.find(inventory.dataMasterInventory(), function(o) {
        return o._id == id;
    });
    ko.mapping.fromJS(data, inventory.record)
    $("#unitDropdown").data('kendoDropDownList').value(inventory.record.Unit());
    $("#AddNewModal").modal("show")
    inventory.text("Update")
}

inventory.delete = function(id, INVID) {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + INVID + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "#3da09a",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function(res) {
        if (res) {
            var url = "/master/deleteinventory";
            var param = {
                ID: id
            };
            ajaxPost(url, param, function(data) {
                if (data.Message != "OK") {
                    swal({
                        title: 'Warning',
                        text: data.Message,
                        type: 'error',
                        confirmButtonColor: "#3da09a"
                    });
                    model.Processing(false);
                } else {
                    swal({
                        title: 'Success',
                        text: 'Inventory has been deleted!',
                        type: 'success',
                        confirmButtonColor: "#3da09a"
                    });
                    inventory.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}

inventory.unitDropdown = function() {
    $("#unitDropdown").html("")
    var data = []
    ajaxPost("/master/getdataunit", {}, function(res) {
        $("#unitDropdown").kendoDropDownList({
            filter: "contains",
            dataTextField: "UnitName",
            dataValueField: "UnitName",
            dataSource: res.Data,
            optionLabel: 'Select one',
            noDataTemplate: $("#noDataTemplate").html(),
            change: function(e) {
                var dataitem = this.dataItem();
                inventory.valueUnit(dataitem.UnitName)
            },
        });
        //console.log(res.Data)
        var dropdownlist = $("#unitDropdown").data("kendoDropDownList");
        if (dropdownlist.value() != "") {
            dropdownlist.value("");
            dropdownlist.trigger("change");
        }
    })
}

inventory.addNewUnit = function(widgetId, value, value2) {
    var widget = $("#" + widgetId).getKendoDropDownList();
    //console.log(widgetId, value, value2, widget)
    var dataSource = widget.dataSource;
    swal({
        title: "Are you sure?",
        text: "You want add new Unit",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: false,
        closeOnCancel: false
    }, function(isConfirm) {
        if (isConfirm) {
            dataSource.add({
                Id: "",
                UnitCode: "",
                UnitName: value.toUpperCase()
            });
            dataSource.one("sync", function() {
                widget.select(dataSource.view().length - 1);
            });
            dataSource.sync();
            var dropdownlist = $("#unitDropdown").data("kendoDropDownList");
            dropdownlist.value(value.toUpperCase());
            dropdownlist.trigger("change");
            model.Processing(true)
            var url = "/master/savenewunit"
            var param = {
                UnitName: value.toUpperCase()
            }
            ajaxPost(url, param, function(data) {
                model.Processing(false)
                swal({
                    title: "Success!",
                    text: "Success Add New Data",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                })
                $('#AddNewModal').modal('hide');
            })
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}

inventory.importExcel = function() {
    event.stopPropagation();
    event.preventDefault();

    model.Processing(true);

    if ($('input[type=file]')[0].files[0] == undefined) {
        swal('Error', 'Please select a file to upload!', 'error');
        model.Processing(false);
        return;
    }
    var len = $("#fDok")[0].files.length
    if (len > 0) {
        var j = 0;
        var data = new FormData();
        for (i = 0; i < $("#fDok")[0].files.length; i++) {
            data.append("filedoc", $('input[type=file]')[0].files[i]);
            data.append("filename", $('input[type=file]')[0].files[i].name);
        }
    } else {
        var data = new FormData();
        data.append("filedoc", $('input[type=file]')[0].files[0]);
        data.append("filename", $('input[type=file]')[0].files[0].name);
    }
    if ($('input[type=file]')[0].files[0].name != "") {
        jQuery.ajax({
            url: '/master/uploadfilesinventory',
            data: data,
            cache: false,
            contentType: false,
            processData: false,
            type: 'POST',
            success: function(data) {
                if (data.IsError) {
                    swal('Error Import Data', data.Message, 'error');
                } else {
                    swal('Success', 'Data has been uploaded successfully!', 'success');
                }
                
                model.Processing(false);
                $('#fDok').val('');
                $('#fInfo').val('');
                $('#ImportModalInventory').modal('hide');
                // inventory.GetDataCoa()
            }
        });
    } else {
        swal('Error', 'Please select a file to upload!', 'error');
        model.Processing(false);
    }
}

inventory.exportExcel = function () {
    $("#gridinventory").getKendoGrid().saveAsExcel();
}

inventory.getDataTypeStock = function () {
    model.Processing(true)
    ajaxPost('/master/getdatatypestock', {}, function (res) {
        if (res.Total === 0) {
            swal({
                title: "Error!", 
                text: res.Message, 
                type: "error",
                confirmButtonColor: "#3da09a"})
            return
        }
        //console.log(res)
        inventory.DataTypeStock(res.Data)  

        model.Processing(false)
    })
}


inventory.searchData = function () {        
    inventory.getDataInventory(function() {
        inventory.renderGrid()
    })  
}

inventory.search = function () {
    var search = inventory.textSearch()
    inventory.renderGrid(search)    
}
inventory.init = function() {    
    inventory.getDataTypeStock()
    inventory.getDataLocation(function() {        
        inventory.getDataInventory(function() {
            inventory.renderGrid()
        })
    })
    inventory.unitDropdown()
    inventory.getDateNow()
    setTimeout(function(){  $("#locationname").data('kendoDropDownList').value(parseInt(userinfo.locationid()));
    ; }, 1000);
   
}

inventory.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["INVDesc", "INVID", "Type", "Unit"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridinventory").data("kendoGrid").dataSource.query({ filter: filter, pageSize: 25, page: 1 });
}

$(function() {
    inventory.init()
    $("#textSearch").on("keyup blur change", function () {
        inventory.filterText();
    });
})