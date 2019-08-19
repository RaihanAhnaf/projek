var transferreceipt = {}
transferreceipt.DatePageBar = ko.observable()
transferreceipt.dataShipment = ko.observableArray([])
transferreceipt.dataMasterInventory = ko.observableArray([])
transferreceipt.dataDropDownInventory = ko.observableArray([])
transferreceipt.gridDataSource = ko.observableArray([])
transferreceipt.datatransferreceipt = ko.observableArray([])
transferreceipt.dataTransferReceiptForReport = ko.observableArray([])
transferreceipt.dataLocation = ko.observableArray([])
transferreceipt.warehouseAll = ko.observableArray([])
transferreceipt.editing = ko.observable(false)
transferreceipt.gridTS = ko.observableArray([])

transferreceipt.newRecord = function () {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumberShipment: "",
        DocumentNumberReceipt: "",
        StoreHouseFrom: "",
        StoreHouseTo: "",
        StoreHouseNameFrom: "",
        StoreHouseNameTo: "",
        Description: "",
        Status: "",
        ListDetailTransferReceipt: [],
    }
    page.ListDetailTransferReceipt.push(transferreceipt.detailTransferReceipt({}))
    return page
}

transferreceipt.detailTransferReceipt = function (data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.Item = data.Item == undefined ? '' : data.Item
    dataTmp.StockUnit = data.StockUnit == undefined ? '' : data.StockUnit
    dataTmp.Qty = data.Qty == undefined ? '' : data.Qty
    dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
    var x = ko.mapping.fromJS(dataTmp)

    return x
}
transferreceipt.record = ko.mapping.fromJS(transferreceipt.newRecord())
transferreceipt.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    transferreceipt.DatePageBar(page)
}

transferreceipt.createTR = function (id) {
    transferreceipt.resetData();
    transferreceipt.editing(true);
    transferreceipt.record.DocumentNumberShipment(id);
    transferreceipt.onShipmentSelected(id);
}

transferreceipt.getDataTransferShipment = function (callback) {
    var startdate = $('#dateStartTS').data('kendoDatePicker').value();
    var enddate = $('#dateEndTS').data('kendoDatePicker').value();

    var param = {
        DateStart: moment(startdate).format('YYYY-MM-DD'),
        DateEnd: moment(enddate).format('YYYY-MM-DD'),
        Filter: true
    }
    model.Processing(true)
    ajaxPost('/transferorder/gettransfershipmentts', param, function (res) {
        if (res.IsError) {
            swal({
                title: "Search Not Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            return
        };
        transferreceipt.dataShipment.removeAll();
        transferreceipt.gridTS.removeAll();
        $(res.Data).each(function (idx, ele) {
            transferreceipt.dataShipment.push({
                text: ele.DocumentNumberShipment,
                value: ele.DocumentNumberShipment,
                obj: ele
            });
            ele.StoreHouseNameFrom = transferreceipt.getStoreHouseName(ele.StoreHouseFrom);
            ele.StoreHouseNameTo = transferreceipt.getStoreHouseName(ele.StoreHouseTo);
            transferreceipt.gridTS.push(ele);
        });
        model.Processing(false)
        if (typeof callback == "function") callback();
    })
}
transferreceipt.refreshDataSource = function (filter, callback) {
    filter = filter || { Filter: false };
    model.Processing(true)
    ajaxPost('/transferorder/gettransferreceipt', filter, function (res) {
        if (res.IsError) {
            model.Processing(false)
            swal({
                title: "Search Not Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            return
        }

        for (var i in res.Data) {
            var data = res.Data[i];
            data.StoreHouseNameFrom = transferreceipt.getStoreHouseName(data.StoreHouseFrom);
            data.StoreHouseNameTo = transferreceipt.getStoreHouseName(data.StoreHouseTo);
        }
        transferreceipt.datatransferreceipt(res.Data)
        transferreceipt.dataTransferReceiptForReport(res.Data)
        transferreceipt.gridDataSource.removeAll();
        $(res.Data).each(function (idx, ele) {
            ele.ListDetailOrder = ele.ListDetailTransferReceipt;
            delete ele.ListDetailTransferReceipt;
            transferreceipt.gridDataSource.push(ele);
        });

        if (filter.Filter) {
            var grid = $("#gridListtransferreceipt").data("kendoGrid");
            if (grid) {
                grid.dataSource.read()
                grid.refresh();
            }
            //console.log("Grid Refreshed: ", res.Data);
        }

        model.Processing(false);
        if (callback) callback.apply(transferreceipt);
        transferreceipt.fromTransferOrderReport()
    }, function () {
        model.Processing(false)
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

transferreceipt.renderGrid = function (filter) {
    this.refreshDataSource(filter, function () {
        if (typeof $('#gridListtransferreceipt').data('kendoGrid') !== 'undefined') {
            $('#gridListtransferreceipt').data('kendoGrid').setDataSource(new kendo.data.DataSource({
                data: transferreceipt.gridDataSource(),
            }))
            return
        }

        var columns = [
            {
                title: 'Action',
                width: 100,
                template: "# if (userinfo.rolenameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' )" +
                    "{#<button onclick='transferreceipt.viewData(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>" +
                    "<button style='display:none' onclick='transferreceipt.editData(\"#: ID #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button>#}" +
                    "else" +
                    "{#<button onclick='transferreceipt.viewData(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>#}#",
            }, {
                field: 'DocumentNumberReceipt',
                title: 'Document Number Receipt',
                width: 160,
            }, {
                field: 'StoreHouseNameFrom',
                title: 'Store House From',
                width: 200,
            }, {
                field: 'StoreHouseNameTo',
                title: 'Store House To',
                width: 200,
            }, {
                field: 'DateStr',
                title: 'Date',
                width: 160,
            }, {
                field: 'Description',
                title: 'Description',
                width: 200,
            }
        ]

        //console.log(transferreceipt.gridDataSource());
        $('#gridListtransferreceipt').kendoGrid({
            dataSource: {
                data: transferreceipt.gridDataSource(),
                sort: {
                    field: 'DateStr',
                    dir: 'desc',
                }
            },
            height: 500,
            width: 140,
            sortable: true,
            scrollable: true,
            columns: columns,
            excelExport: function (e) {
                ProActive.kendoExcelRender(e, "TransferReceipt", function (row, sheet) {
                    for (var ci = 0; ci < row.cells.length; ci++) {
                        var cell = row.cells[ci];
                        if (row.type == "data") {
                            if (ci == 1 || ci == 2)
                                cell.hAlign = "left";
                            if (ci == 3)
                                cell.format = "dd-MM-yyyy";
                        }
                    }
                });
            },
        });
    });
};
transferreceipt.exportExcel = function () {
    $("#gridListtransferreceipt").getKendoGrid().saveAsExcel();
}

transferreceipt.addNewItem = function () {
    transferreceipt.record.ListDetailTransferReceipt.push(transferreceipt.detailTransferReceipt({}));
    $(".invhide").show();
}
transferreceipt.getStoreHouseName = function (id) {
    var wh = transferreceipt.warehouseAll();
    for (var i in wh) {
        var w = wh[i];
        if (w.value == id)
            return w.value + " - " + w.text;
    }
    return id;
},
    transferreceipt.removeRow = function () {
        transferreceipt.record.ListDetailTransferReceipt.remove(this)
        if (transferreceipt.record.ListDetailTransferReceipt().length == 0) {
            transferreceipt.record.ListDetailTransferReceipt.push(transferreceipt.detailTransferReceipt({}))
        }
    }
transferreceipt.resetData = function () {
    var rec = transferreceipt.newRecord();
    rec.DatePosting = moment().format("DD-MMM-YYYY");
    ko.mapping.fromJS(rec, transferreceipt.record);
    $("#btnSave").prop("disabled", false);
    $("#btnDelete").prop("disabled", false);
    $("#btnSave").prop("disabled", false);
    $("#btnDelete").hide();
    $("#btnPrint").hide();
    $("#btnReset").show();
    $("#btnSave").show();
    $("#storehousefrom").data('kendoDropDownList').enable(true);
    $("#datepurchase").data('kendoDatePicker').enable(true);
    $("#taDesc").prop("disabled", false);
    $(".hide-on-view").hide();
    transferreceipt.editing(false);
    transferreceipt.renderGridTS();
}
transferreceipt.renderGridTS = function () {
    this.getDataTransferShipment(function () {
        $('#gridTS').html('');

        var columns = [
            {
                title: 'Action',
                width: 100,
                template: "<button onclick='transferreceipt.createTR(\"#: DocumentNumberShipment #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil' aria-hidden='true'></i></button>",
            }, {
                field: 'DocumentNumberShipment',
                title: 'Document Number Shipment',
                width: 160,
            }, {
                field: 'StoreHouseNameFrom',
                title: 'Store House From',
                width: 200,
            }, {
                field: 'StoreHouseNameTo',
                title: 'Store House To',
                width: 200,
            }, {
                field: 'DateStr',
                title: 'Date',
                width: 160,
            }, {
                field: 'Description',
                title: 'Description',
                width: 200,
            }
        ]

        //console.log(transferreceipt.gridDataSource());
        $('#gridTS').kendoGrid({
            dataSource: {
                data: transferreceipt.gridTS(),
                sort: {
                    field: 'DateStr',
                    dir: 'desc',
                }
            },
            height: 500,
            width: 140,
            sortable: true,
            scrollable: true,
            columns: columns,
            excelExport: function (e) {
                ProActive.kendoExcelRender(e, "TransferReceipt", function (row, sheet) {
                    for (var ci = 0; ci < row.cells.length; ci++) {
                        var cell = row.cells[ci];
                        if (row.type == "data") {
                            if (ci == 1 || ci == 2)
                                cell.hAlign = "left";
                            if (ci == 3)
                                cell.format = "dd-MM-yyyy";
                        }
                    }
                });
            },
        });
    });
}
transferreceipt.onChangeCodeItem = function (value, idx) {
    var result = _.filter(transferreceipt.dataDropDownInventory(), {
        'Kode': value
    })[0];
    transferreceipt.record.ListDetailTransferReceipt()[idx].Item(result ? result.Name : undefined);
    //console.log(result);
    transferreceipt.record.ListDetailTransferReceipt()[idx].StockUnit(result ? result.Saldo : undefined);
}

transferreceipt.deleteData = function () {
    //console.log(transferreceipt.record);
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + transferreceipt.record.DocumentNumberReceipt() + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            var url = "/transferorder/deletetransferreceipt";
            var param = {
                ID: transferreceipt.record.ID()
            };
            ajaxPost(url, param, function (e) {
                if (e.Message == "OK") {
                    setTimeout(function () {
                        swal({
                            title: "Success!",
                            text: "Data has been deleted!",
                            type: "success",
                            confirmButtonColor: "#3da09a"
                        }, function () {
                            window.location.assign("/transferorder/transferreceipt")
                        });
                    }, 100)
                } else {
                    swal({
                        title: 'Warning',
                        text: e,
                        type: 'error',
                        confirmButtonColor: "#3da09a"
                    });
                    model.Processing(false);
                }
            }, undefined);
        } else {
            swal({
                title: "Cancelled",
                text: "",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
            model.Processing(false);
        }
    });
}
transferreceipt.search = function () {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {
        DateStart: moment(startdate).format('YYYY-MM-DD'),
        DateEnd: moment(enddate).format('YYYY-MM-DD'),
        Filter: true
    }
    transferreceipt.refreshDataSource(param);
}
transferreceipt.saveData = function () {
    var change = ko.mapping.toJS(transferreceipt.record);

    if (change.DocumentNumberShipment == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't selected a Document Number Shipment",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }

    change.StoreHouseFrom = parseInt(change.StoreHouseFrom)
    change.StoreHouseTo = parseInt(change.StoreHouseTo)
    for (var i = 0; i < change.ListDetailTransferReceipt.length; i++) {
        change.ListDetailTransferReceipt[i].Qty = parseInt(change.ListDetailTransferReceipt[i].Qty)
        change.ListDetailTransferReceipt[i].StockUnit = parseInt(change.ListDetailTransferReceipt[i].StockUnit)

        if (change.ListDetailTransferReceipt[i].Qty <= 0) {
            return swal({
                title: 'Warning!',
                text: "Item: " + change.ListDetailOrder[i].Item + " has invalid Qty",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        }
        if (change.ListDetailTransferReceipt[i].Item == "") {
            return swal({
                title: 'Warning!',
                text: "Cannot save blank item",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        }
    };
    change.DatePosting = new Date($('#datepurchase').data('kendoDatePicker').value());
    (change.DatePosting).setHours((change.DatePosting).getHours() + 7)
    change.DateStr = moment($('#datepurchase').data('kendoDatePicker').value()).format("DD-MMM-YYYY");

    if (change.DateStr == "Invalid date") {
        return swal({
            title: 'Warning!',
            text: "You haven't entered a date",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    change.Status = "TS";

    var param = {
        Data: change
    };

    var url = "/transferorder/savetransferreceipt"
    swal({
        title: "Are you sure?",
        text: "You will submit this Transfer Receipt",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, param, function (e) {
                setTimeout(function () {
                    swal({
                        title: "Success!",
                        text: "Data has been saved!",
                        type: "success",
                        confirmButtonColor: "#3da09a"
                    }, function () {
                        // window.location.assign("/transferorder/transferreceipt")
                        $("#btnSave").prop('disabled', false);
                        transferreceipt.resetData();
                        transferreceipt.createdForm();
                    });
                }, 100)
                model.Processing(false)
            })
        } else {
            swal({
                title: "Cancelled",
                text: "",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}
transferreceipt.getMasterInventory = function (fromLocation, callback) {
    if (fromLocation) {
        var param = {
            Filter: true,
            LocationID: parseInt(fromLocation)
        };
        ajaxPost("/transferorder/getmasterinventory", param, function (json) {
            if (!json.IsError) {
                //console.log(json);
                transferreceipt.dataDropDownInventory.removeAll();
                $(json.Data).each(function (ix, ele) {
                    //console.log(ele);
                    transferreceipt.dataDropDownInventory.push({
                        value: ele.INVID,
                        text: ele.INVID,
                        name: ele.INVDesc,
                        stock: ele.Saldo
                    });
                });

            }
            if (callback) callback();
        });
    }
    else {
        transferreceipt.dataDropDownInventory.removeAll();
    }
}
transferreceipt.onShipmentSelected = function (value) {
    var data = value;
    var result = _.filter(transferreceipt.dataShipment(), {
        'value': data
    })[0];
    if (result.obj) {
        var rec = {};
        Object.assign(rec, result.obj);
        rec.ID = "";
        rec.Description = "";
        rec.DateStr = moment().format("DD-MMM-YYYY");
        rec.DatePosting = moment().format("DD-MMM-YYYY");
        rec.ListDetailTransferReceipt = [];
        for (var k in rec.ListDetailTransferShipment) {
            rec.ListDetailTransferReceipt.push(rec.ListDetailTransferShipment[k]);
        }
        delete rec.ListDetailTransferShipment;
        transferreceipt.getMasterInventory(rec.StoreHouseFrom, function () {
            ko.mapping.fromJS(rec, transferreceipt.record);
        });

        $("#tsStoreHouseFrom").val(result.obj.StoreHouseNameFrom);
        $("#tsStoreHouseTo").val(result.obj.StoreHouseNameTo);
    }
}

transferreceipt.createdForm = function () {
    transferreceipt.resetData();
},

    transferreceipt.backToList = function () {
        $('.nav-tabs a[href="#List"]').tab('show');

        //Remove Parameter that exist in URL
        var url_string = window.location.href
        var url = new URL(url_string);
        var param = url.searchParams.get("id");
        if (param != null) {
            window.history.pushState("", "", "/" + "transferorder/transferreceipt");
            location.reload();
        }

        transferreceipt.refreshDataSource();
    }

transferreceipt.editData = function (e) {
    data = _.find(transferreceipt.gridDataSource(), function (o) {
        return o.ID == e;
    });
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY');

    // Remove all shipment list while adding edited shipment doc number
    transferreceipt.record.ListDetailTransferReceipt.removeAll();

    transferreceipt.dataShipment.removeAll();
    transferreceipt.dataShipment.push({
        text: data.DocumentNumberShipment,
        value: data.DocumentNumberShipment,
    });
    transferreceipt.record.DocumentNumberShipment(data.DocumentNumberShipment)
    ko.mapping.fromJS(data, transferreceipt.record);
    transferreceipt.getMasterInventory(data.StoreHouseTo, function () {
        transferreceipt.record.ListDetailTransferReceipt.removeAll();
        $(data.ListDetailOrder).each(function (idx, ele) {
            //console.log("Ele#" + idx, ele);
            transferreceipt.record.ListDetailTransferReceipt.push(ele);
        });
    })

    $("#btnSave").prop("disabled", false);
    $("#btnDelete").prop("disabled", false);
    $("#btnSave").show();
    $("#btnDelete").show();
    $("#btnPrint").hide();
    $("#btnReset").hide();
    $('.nav-tabs a[href="#Create"]').tab('show');
    $("#storehousefrom").data('kendoDropDownList').enable(true);
    $("#datepurchase").data('kendoDatePicker').enable(true);
    $("#taDesc").prop("disabled", false);
    $("#buttonAdd").prop("disabled", true);
    $(".hide-on-view").hide();
}
transferreceipt.viewData = function (e) {
    transferreceipt.editData(e);
    transferreceipt.editing(true);

    $("#btnSave").prop("disabled", true);
    $("#btnDelete").prop("disabled", true);
    $("#taDesc").prop("disabled", true);
    $("#btnSave").hide();
    $("#btnDelete").hide();
    $("#btnReset").hide();
    $("#btnPrint").show();
    $("#storehousefrom").data('kendoDropDownList').enable(false);
    $("#datepurchase").data('kendoDatePicker').enable(false);
    setTimeout(() => {
        $(".hide-on-view").hide();
    }, 100);
}

// Print 
transferreceipt.getMasterLocation = function (callbackall) {
    $.ajax({
        url: "/transferorder/getdatamasterlocation",
        success: function (json) {
            if (!json.IsError) {
                $(json.Data).each(function (ix, ele) {
                    transferreceipt.dataLocation.push({
                        value: ele.LocationID,
                        text: ele.LocationName
                    });
                });
            }
        }
    });
    $.ajax({
        url: "/transferorder/getalllocations",
        success: function (json) {
            transferreceipt.warehouseAll.removeAll();
            $(json).each(function (ix, ele) {
                transferreceipt.warehouseAll.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
                if (typeof callbackall == "function") {
                    callbackall();
                }
            });
        }
    });
}
transferreceipt.newRecordPrint = function () {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumberShipment: "",
        DocumentNumberReceipt: "",
        StoreHouseFrom: "",
        StoreHouseTo: "",
        Description: "",
        Status: "",
        ListDetailTransferReceipt: []
    }
    page.ListDetailTransferReceipt.push(transferreceipt.detailPrintTransferReceipt({}))
    return page
}

transferreceipt.detailPrintTransferReceipt = function (data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.Item = data.Item == undefined ? '' : data.Item
    dataTmp.StockUnit = data.StockUnit == undefined ? '' : data.StockUnit
    dataTmp.Qty = data.Qty == undefined ? '' : data.Qty
    dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
    var x = ko.mapping.fromJS(dataTmp)
    return x
}

transferreceipt.recordPrint = ko.mapping.fromJS(transferreceipt.newRecordPrint());

transferreceipt.print = function () {
    $("#printtransferreceipt").show();
    $("#transferreceipt").hide();
    var allData = {}
    var e = transferreceipt.record.ID()
    allData = transferreceipt.datatransferreceipt()
    var data = _.find(allData, function (o) {
        return o.ID == e;
    });
    //console.log("data")
    //console.log(data)
    var dataFrom = _.find(transferreceipt.dataLocation(), function (n) {
        return n.value == data.StoreHouseFrom
    });
    //console.log("dataFrom")
    //console.log(dataFrom)

    var dataTo = _.find(transferreceipt.dataLocation(), function (b) {
        return b.value == data.StoreHouseTo
    });

    transferreceipt.recordPrint.ID(data.ID)
    transferreceipt.recordPrint.DateStr(data.DateStr)
    transferreceipt.recordPrint.DocumentNumberReceipt(data.DocumentNumberReceipt)
    transferreceipt.recordPrint.StoreHouseFrom(dataFrom.text)
    transferreceipt.recordPrint.StoreHouseTo(dataTo.text)
    transferreceipt.recordPrint.Description(data.Description)
    transferreceipt.recordPrint.Status(data.Status)
    transferreceipt.recordPrint.ListDetailTransferReceipt([])
    for (i in data.ListDetailOrder) {
        transferreceipt.recordPrint.ListDetailTransferReceipt.push(transferreceipt.detailPrintTransferReceipt({}))
        transferreceipt.recordPrint.ListDetailTransferReceipt()[i].CodeItem(data.ListDetailOrder[i].CodeItem)
        transferreceipt.recordPrint.ListDetailTransferReceipt()[i].Item(data.ListDetailOrder[i].Item)
        transferreceipt.recordPrint.ListDetailTransferReceipt()[i].Qty(data.ListDetailOrder[i].Qty)
    }
    window.print();
    transferreceipt.reload()
}
transferreceipt.printPdf = function () {
    var param = {
        Id: transferreceipt.record.ID()
    }
    model.Processing(true)
    ajaxPost("/transferorder/ExportPdfPerDataTR".toLowerCase(), param, function (res) {
        if (res.IsError) {
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}
transferreceipt.reload = function () {
    $("#printtransferreceipt").hide();
    $("#transferreceipt").show();
}
transferreceipt.textSearch = function () { }
transferreceipt.initDate = function () {
    ProActive.KendoDatePickerRange();
    ProActive.KendoDatePickerRange("dateStartTS", "dateEndTS");
}
transferreceipt.init = function () {
    transferreceipt.initDate();
    transferreceipt.getMasterInventory();
    transferreceipt.getMasterLocation(function () {
        transferreceipt.getDateNow();
        transferreceipt.getDataTransferShipment();
        var filter = {
            DateStart: moment().format('YYYY-MM-') + "01",
            DateEnd: moment().format('YYYY-MM-DD'),
            Filter: true
        }
        transferreceipt.renderGrid(filter);

        transferreceipt.resetData();
    });
    ProActive.GlobalSearch("gridListtransferreceipt",
        [
            "DocumentNumberShipment",
            "DocumentNumberReceipt",
            "Description",
            "StoreHouseNameFrom",
            "StoreHouseNameTo"
        ]
    );
}
//===End Print==========

transferreceipt.fromTransferOrderReport = function () {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null) {
        var allData = transferreceipt.dataTransferReceiptForReport()

        var data = _.find(allData, function (o) {
            return o.DocumentNumberReceipt == num;
        });
        //  console.log(purchaseorder.dataMasterPurchaseInventory())
        if (data != undefined) {
            transferreceipt.viewData(data.ID)
            setTimeout(() => {
                transferreceipt.dataShipment.push({
                    text: data.DocumentNumberShipment,
                    value: data.DocumentNumberShipment,
                });
                transferreceipt.record.DocumentNumberShipment.valueHasMutated();
            }, 100);
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning",
                confirmButtonColor: "#3da09a"
            }, function () {
                window.location.assign("/transferorder/transferreceipt")
            });
        }
    }
}

$(function () {
    transferreceipt.init()
});