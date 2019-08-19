model.Processing(false);

// TransferShipment Object BEGIN ================
var transferShipment = {
    warehouseFrom: ko.observableArray([]),
    warehouseTo: ko.observableArray([]),
    warehouseAll: ko.observableArray([]),
    gridDataSource: ko.observableArray([]),
    dataTransferShipment: ko.observableArray([]),
    dataTransferShipmentForReport: ko.observableArray([]),
    dataLocation: ko.observableArray([]),
    filterindicator: ko.observable(false),
    DatePageBar: ko.observable(),
    dataDropDownInventory: ko.observableArray([]),
    storehouseData : [],


    getDateNow: function () {
        var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
        this.DatePageBar(page)
    },

    newRecord: function () {
        var page = {
            ID: "",
            DateStr: "",
            DatePosting: "",
            DocumentNumber: "",
            StoreHouseFrom: "",
            StoreHouseTo: "",
            StoreHouseNameFrom: "",
            StoreHouseNameTo: "",
            Description: "",
            ListDetailOrder: [],
        }
        page.ListDetailOrder.push(this.listDetailOrder({}))
        return page
    },

    newListDetailOrder: function () {
        return {
            Id: "",
            CodeItem: "",
            Item: "",
            StockUnit: undefined,
            Qty: undefined
        };
    },

    listDetailOrder: function (data) {
        var dataTmp = {}
        dataTmp.Id = data.Id == undefined ? '' : data.Id
        dataTmp.Item = data.Item == undefined ? '' : data.Item
        dataTmp.StockUnit = data.StockUnit == undefined ? '' : data.StockUnit
        dataTmp.Qty = data.Qty == undefined ? '' : data.Qty
        dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
        var x = ko.mapping.fromJS(dataTmp)
        return x
    },

    createdForm: function () {
        transferShipment.resetData();
    },

    getMasterLocation: function (callbackall) {
        $.ajax({
            url: "/transferorder/getuserlocations",
            success: function (json) {
                transferShipment.warehouseFrom.removeAll();
                $(json).each(function (ix, ele) {
                    transferShipment.dataLocation.push({
                        value: ele.LocationID,
                        text: ele.LocationName
                    });
                    transferShipment.warehouseFrom.push({
                        value: ele.LocationID,
                        text: ele.LocationName
                    });
                });
            }
        });
        $.ajax({
            url: "/transferorder/getuserdestlocations",
            success: function (json) {
                transferShipment.warehouseTo.removeAll();
                $(json).each(function (ix, ele) {
                    transferShipment.warehouseTo.push({
                        value: ele.LocationID,
                        text: ele.LocationName
                    });
                });
            }
        });
        $.ajax({
            url: "/transferorder/getalllocations",
            success: function (json) {
                transferShipment.storehouseData = json;
                transferShipment.warehouseAll.removeAll();
                $(json).each(function (ix, ele) {
                    transferShipment.warehouseAll.push({
                        value: ele.LocationID,
                        text: ele.LocationName
                    });
                });
                if (typeof callbackall == "function")
                    callbackall.apply(transferShipment);
            }
        });
    },

    getMasterInventory: function (fromLocation, callback) {
        if (fromLocation) {
            var param = {
                Filter: true,
                LocationID: parseInt(fromLocation)
            };
            ajaxPost("/transferorder/getmasterinventory", param, function (json) {
                if (!json.IsError) {
                    //console.log(json);
                    transferShipment.dataDropDownInventory.removeAll();
                    $(json.Data).each(function (ix, ele) {
                        transferShipment.dataDropDownInventory.push({
                            value: ele.INVID,
                            text: ele.INVID,
                            name: ele.INVDesc,
                            stock: ele.Saldo
                        });
                    });

                }
                if (callback) callback();
            });
        } else {
            transferShipment.dataDropDownInventory.removeAll();
            if (callback) callback();
        }
    },

    init: function () {
        this.initDate();
        this.getDateNow();
        this.getMasterLocation(function () {
            this.getMasterInventory();

            var filter = {
                DateStart: moment().format('YYYY-MM-') + "01",
                DateEnd: moment().format('YYYY-MM-DD'),
                Filter: true
            }
            this.renderGrid(filter);
            this.resetData();
        });
        
        ProActive.GlobalSearch("gridListTransferShipment",
            [
                "DocumentNumberShipment",
                "DocumentNumberReceipt",
                "Description",
                "StoreHouseNameFrom",
                "StoreHouseNameTo",
                "StoreHouseFrom equals",
                "StoreHouseTo equals"
            ]
        );
    },

    getStoreHouseName: function (id) {
        var wh = transferShipment.warehouseAll();
        for (var i in wh) {
            var w = wh[i];
            if (w.value == id)
                return w.value + " - " + w.text;
        }
        return id;
    },

    refreshDataSource: function (filter, callback) {
        filter = filter || {
            Filter: false
        };
        model.Processing(true)
        ajaxPost('/transferorder/gettransfershipment', filter, function (res) {
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
            model.Processing(false)
            for (var i in res.Data) {
                var data = res.Data[i];
                data.StoreHouseNameFrom = transferShipment.getStoreHouseName(data.StoreHouseFrom);
                data.StoreHouseNameTo = transferShipment.getStoreHouseName(data.StoreHouseTo);
            }
            transferShipment.dataTransferShipment(res.Data)
            transferShipment.dataTransferShipmentForReport(res.Data)
            transferShipment.gridDataSource.removeAll();
            $(res.Data).each(function (idx, ele) {
                ele.ListDetailOrder = ele.ListDetailTransferShipment;
                delete ele.ListDetailTransferShipment;
                transferShipment.gridDataSource.push(ele);
            });

            if (filter.Filter) {
                var grid = $("#gridListTransferShipment").data("kendoGrid");
                if (grid) {
                    grid.dataSource.read()
                    grid.refresh();
                }
                //console.log("Grid Refreshed: ", res.Data);
            }

            model.Processing(false);
            if (callback) callback.apply(transferShipment);
            transferShipment.fromTransferOrderReport();
        }, function () {
            model.Processing(false)
            swal({
                title: "Error!",
                text: "Unknown error, please try again",
                type: "error",
                confirmButtonColor: "#3da09a"
            })
        })
    },

    search: function () {
        var startdate = $('#dateStart').data('kendoDatePicker').value();
        var enddate = $('#dateEnd').data('kendoDatePicker').value();

        var param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true
        }
        transferShipment.refreshDataSource(param);
    },

    renderGrid: function (filter) {
        this.refreshDataSource(filter, function () {
            if (typeof $('#gridListTransferShipment').data('kendoGrid') !== 'undefined') {
                $('#gridListTransferShipment').data('kendoGrid').setDataSource(new kendo.data.DataSource({
                    data: transferShipment.gridDataSource(),
                }))
                return
            }

            var columns = [{
                title: 'Action',
                width: 100,
                template: "# if (userinfo.rolenameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' )" +
                    "{#<button onclick='transferShipment.viewData(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>" +
                    "<button style='display:none' onclick='transferShipment.editData(\"#: ID #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button>#}" +
                    "else" +
                    "{#<button onclick='transferShipment.viewData(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>#}#",
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
            }]

            $('#gridListTransferShipment').kendoGrid({
                dataSource: {
                    data: transferShipment.gridDataSource(),
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
                    ProActive.kendoExcelRender(e, "TransferShipment", function (row, sheet) {
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
    },
    exportExcel: function () {
        $("#gridListTransferShipment").getKendoGrid().saveAsExcel();
    },

    listDetail: function (data) {
        var dataTmp = this.newListDetailOrder();
        Object.assign(dataTmp, data);

        return ko.mapping.fromJS(dataTmp)
    },

    addNewItem: function () {
        transferShipment.record.ListDetailOrder.push(transferShipment.listDetail({}));
        //this.checkdata()
        //this.maskingMoney()
        $(".invhide").show();
    },

    getStoreHouse: function(locID, relation) {
        locID = parseInt(locID);
        relation = relation || "self";
        var self = null;
        var parent = null;
        var children = [];
        for (var k in transferShipment.storehouseData) {
            var sh = transferShipment.storehouseData[k];
            if (sh.LocationID == locID) {
                self = sh;
                if (relation == "self")
                    return sh;
            }
            if (sh.Main_LocationID == locID && (relation == "children" || relation == "relatives")) {
                children.push(sh);
            }
        }
        if ((relation == "parent" || relation == "relatives") && self != null) {
            for (var k in transferShipment.storehouseData) {
                var sh = transferShipment.storehouseData[k];
                if (sh.LocationID == self.Main_LocationID) {
                    parent = sh;
                    if (relation == "parent")
                        return parent;
                    break;
                }
            }
        }
        if (relation == "relatives" && parent != null) {
            children = [parent].concat(children);
        }
        return relation == "self" ? null : children;
    },

    fillStoreHouseTo: function(fromLocation) {
        var locs = transferShipment.getStoreHouse(fromLocation, "relatives");
        transferShipment.warehouseTo.removeAll();
        for(var k in locs)
        {
            transferShipment.warehouseTo.push({
                value: locs[k].LocationID,
                text: locs[k].LocationName
            });
        }
    },

    onChangeStoreHouseFrom: function (data, callback) {
        transferShipment.getMasterInventory(data, function(){
            transferShipment.fillStoreHouseTo(data);
            if (typeof callback == "function") callback();
            $("#storehouseto").data('kendoDropDownList').enable(!$("#btnSave").prop("disabled") && data !== "");
            transferShipment.record.StoreHouseTo.valueHasMutated();
        });        
    },

    onChangeCodeItem: function (value, idx) {
        var result = _.filter(transferShipment.dataDropDownInventory(), {
            'value': value
        })[0];
        transferShipment.record.ListDetailOrder()[idx].Item(result ? result.name : undefined);
        transferShipment.record.ListDetailOrder()[idx].StockUnit(result ? result.stock : undefined);
    },

    removeRow: function () {
        transferShipment.record.ListDetailOrder.remove(this)
        if (transferShipment.record.ListDetailOrder().length == 0) {
            transferShipment.record.ListDetailOrder.push(transferShipment.listDetail({}))
        }
    },

    saveData: function () {
        var change = ko.mapping.toJS(transferShipment.record);
        if (change.StoreHouseFrom == "") {
            return swal({
                title: 'Warning!',
                text: "You haven't choose the Store House From",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        };
        if (change.StoreHouseTo == "") {
            return swal({
                title: 'Warning!',
                text: "You haven't choose the Store House To",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        };
        if (change.StoreHouseTo == change.StoreHouseFrom) {
            return swal({
                title: 'Warning!',
                text: "Cannot move to the same Store House",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        };
        change.StoreHouseFrom = parseInt(change.StoreHouseFrom)
        change.StoreHouseTo = parseInt(change.StoreHouseTo)
        for (var i = 0; i < change.ListDetailOrder.length; i++) {
            change.ListDetailOrder[i].Qty = parseInt(change.ListDetailOrder[i].Qty)
            change.ListDetailOrder[i].StockUnit = parseInt(change.ListDetailOrder[i].StockUnit)

            if (change.ListDetailOrder[i].Item == "") {
                return swal({
                    title: 'Warning!',
                    text: "Cannot save blank item",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            }
            if (!change.ListDetailOrder[i].Qty || change.ListDetailOrder[i].Qty <= 0) {
                return swal({
                    title: 'Warning!',
                    text: "Item: " + change.ListDetailOrder[i].Item + " has invalid Qty",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            }
            if (change.ListDetailOrder[i].Qty > change.ListDetailOrder[i].StockUnit) {
                return swal({
                    title: 'Warning!',
                    text: "Cannot move " + change.ListDetailOrder[i].Item + " because Unit Stock is lower than Qty",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            }
        };
        change.ListDetailTransferShipment = change.ListDetailOrder;
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

        var url = "/transferorder/savetransfershipment"
        swal({
            title: "Are you sure?",
            text: "You will submit this Transfer Shipment",
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
                            //window.location.assign("/transferorder/transfershipment")
                            transferShipment.resetData();
                            $("#btnSave").prop("disabled", false)
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
    },

    resetData: function () {
        transferShipment.getMasterLocation();
        var rec = transferShipment.newRecord();
        rec.DatePosting = moment().format("DD-MMM-YYYY");
        ko.mapping.fromJS(rec, transferShipment.record);
        transferShipment.onChangeStoreHouseFrom();

        $("#btnSave").prop("disabled", false);
        $("#btnDelete").prop("disabled", false);
        $("#btnDelete").hide();
        $("#btnPrint").hide();
        $("#btnReset").show();
        $("#btnSave").show();
        $("#storehouseto").data('kendoDropDownList').enable(true);
        $("#storehousefrom").data('kendoDropDownList').enable(true);
        $("#datepurchase").data('kendoDatePicker').enable(true);
        $("#storehouseto").data('kendoDropDownList').enable(false);
        $("select.editableDD").each(function (idx, ele) {
            $(this).data("kendoDropDownList").enable(true);
        });
        $("input.editableInput").each(function (idx, ele) {
            $(this).prop("disabled", false);
        });
        $("#taDesc").prop("disabled", false);
        $(".hide-on-view").show();
    },

    editData: function (e) {
        data = _.find(transferShipment.gridDataSource(), function (o) {
            return o.ID == e;
        });
        $(transferShipment.warehouseAll()).each(function (idx, ele) {
            if (ele.value == data.StoreHouseTo) {
                transferShipment.warehouseTo.push({
                    value: ele.value,
                    text: ele.text
                });
            }
        });
        data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY');
        ko.mapping.fromJS(data, transferShipment.record);
        transferShipment.onChangeStoreHouseFrom(transferShipment.record.StoreHouseFrom(), function () {
            console.log(transferShipment.warehouseTo().length)
            _.each(transferShipment.record.ListDetailOrder(), function (v, i) {
                if (v != undefined) {
                    v.CodeItem.valueHasMutated();
                    v.StockUnit.valueHasMutated();
                }
            });
        });

        if (transferShipment.record.ListDetailOrder().length == 0) {
            transferShipment.record.ListDetailOrder.push(transferShipment.listDetail({}))
        };
        $("#btnSave").prop("disabled", false);
        $("#btnDelete").prop("disabled", false);
        $('.nav-tabs a[href="#Create"]').tab('show');
        $("#btnSave").show();
        $("#btnDelete").show();
        $("#btnPrint").hide();
        $("#btnReset").hide();
        $("#storehouseto").data('kendoDropDownList').enable(true);
        $("#storehousefrom").data('kendoDropDownList').enable(true);
        $("#datepurchase").data('kendoDatePicker').enable(true);
        $("#storehouseto").data('kendoDropDownList').enable(true);
        $("select.editableDD").each(function (idx, ele) {
            $(this).data("kendoDropDownList").enable(true);
        });
        $("input.editableInput").each(function (idx, ele) {
            $(this).prop("disabled", false);
        });
        $("#taDesc").prop("disabled", false);
        $(".hide-on-view").show();
    },

    viewData: function (e) {
        transferShipment.editData(e);

        $("#btnSave").prop("disabled", true);
        $("#btnDelete").prop("disabled", true);
        $("#btnSave").hide();
        $("#btnDelete").hide();
        $("#btnReset").hide();
        $(".hide-on-view").hide();
        $("#btnPrint").show();
        $("#datepurchase").data('kendoDatePicker').enable(false);
        $("#storehouseto").data('kendoDropDownList').enable(false);
        $("#storehousefrom").data('kendoDropDownList').enable(false);
        $("#datepurchase").data('kendoDatePicker').enable(false);
        $("#taDesc").prop("disabled", true);
        $("select.editableDD").each(function (idx, ele) {
            $(this).data("kendoDropDownList").enable(false);
        });
        $("input.editableInput").each(function (idx, ele) {
            $(this).prop("disabled", true);
        });
    },

    deleteData: function () {
        model.Processing(true);
        swal({
            title: "Are you sure to delete " + ko.mapping.toJS(transferShipment.record).DocumentNumber + "?",
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
                var url = "/transferorder/deletetransfershipment";
                var param = {
                    ID: ko.mapping.toJS(transferShipment.record).ID
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
                                window.location.assign("/transferorder/transfershipment")
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
};

transferShipment.record = ko.mapping.fromJS(transferShipment.newRecord());

transferShipment.backToList = function () {
    $('.nav-tabs a[href="#List"]').tab('show');

    //Remove Parameter that exist in URL
    var url_string = window.location.href
    var url = new URL(url_string);
    var param = url.searchParams.get("id");
    if (param != null) {
        window.history.pushState("", "", "/" + "transferorder/transfershipment");
    }
    location.reload();

    transferShipment.refreshDataSource();
}

transferShipment.newRecordPrint = function () {
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
        ListDetailTransferShipment: []
    }
    page.ListDetailTransferShipment.push(transferShipment.newListDetailOrderPrint({}))
    return page
}

transferShipment.newListDetailOrderPrint = function (data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.Item = data.Item == undefined ? '' : data.Item
    dataTmp.StockUnit = data.StockUnit == undefined ? '' : data.StockUnit
    dataTmp.Qty = data.Qty == undefined ? '' : data.Qty
    dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
    var x = ko.mapping.fromJS(dataTmp)
    return x
}

transferShipment.recordPrint = ko.mapping.fromJS(transferShipment.newRecordPrint());

transferShipment.print = function () {
    $("#printtransferShipment").show();
    $("#transferShipment").hide();
    var allData = {}
    var e = transferShipment.record.ID()
    allData = transferShipment.dataTransferShipment()
    var data = _.find(allData, function (o) {
        return o.ID == e;
    });

    var dataFrom = _.find(transferShipment.dataLocation(), function (n) {
        return n.value == data.StoreHouseFrom
    });

    var dataTo = _.find(transferShipment.dataLocation(), function (b) {
        return b.value == data.StoreHouseTo
    });


    transferShipment.recordPrint.ID(data.ID)
    transferShipment.recordPrint.DateStr(data.DateStr)
    transferShipment.recordPrint.DocumentNumberShipment(data.DocumentNumberShipment)
    transferShipment.recordPrint.StoreHouseFrom(dataFrom.text)
    transferShipment.recordPrint.StoreHouseTo(dataTo.text)
    transferShipment.recordPrint.Description(data.Description)
    transferShipment.recordPrint.Status(data.Status)
    transferShipment.recordPrint.ListDetailTransferShipment([])
    for (i in data.ListDetailOrder) {
        transferShipment.recordPrint.ListDetailTransferShipment.push(transferShipment.newListDetailOrderPrint({}))
        transferShipment.recordPrint.ListDetailTransferShipment()[i].CodeItem(data.ListDetailOrder[i].CodeItem)
        transferShipment.recordPrint.ListDetailTransferShipment()[i].Item(data.ListDetailOrder[i].Item)
        transferShipment.recordPrint.ListDetailTransferShipment()[i].Qty(data.ListDetailOrder[i].Qty)
    }
    window.print();
    transferShipment.reload()
}

transferShipment.printPdf = function () {
    var param = {
        Id: transferShipment.record.ID()
    }
    model.Processing(true)
    ajaxPost("/transferorder/ExportPdfPerDataTS".toLowerCase(), param, function (res) {
        if (res.IsError) {
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}

transferShipment.fromTransferOrderReport = function () {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null) {
        var allData = transferShipment.dataTransferShipmentForReport()

        var data = _.find(allData, function (o) {
            return o.DocumentNumberShipment == num;
        });
        //  console.log(purchaseorder.dataMasterPurchaseInventory())
        if (data != undefined) {
            transferShipment.viewData(data.ID)
            setTimeout(() => {
                transferShipment.record.StoreHouseFrom.valueHasMutated();
                transferShipment.record.StoreHouseTo.valueHasMutated();
            }, 100);
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning",
                confirmButtonColor: "#3da09a"
            }, function () {
                window.location.assign("/transferorder/transfershipment")
            });
        }
    }
}

transferShipment.reload = function () {
    $("#printtransferShipment").hide();
    $("#transferShipment").show();
}
transferShipment.initDate = function () {
    var dtpStart = $('#dateStart').data('kendoDatePicker');
    var dtpEnd = $('#dateEnd').data('kendoDatePicker');
    dtpStart.value(moment().startOf('month').toDate());
    dtpEnd.value(moment().startOf('day').toDate());
    dtpStart.max(dtpEnd.value());
    dtpEnd.min(dtpStart.value());

    dtpStart.bind("change", function () {
        dtpEnd.min(dtpStart.value());
    });
    dtpEnd.bind("change", function () {
        dtpStart.max(dtpEnd.value());
    });
}
// TransferShipment Object END ==================


$(function () {
    transferShipment.init();
});
