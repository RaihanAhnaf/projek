var salescreditmemo = {}
var processs = 0;
salescreditmemo.date = ko.observable()
salescreditmemo.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
salescreditmemo.DateEnd = ko.observable(new Date)
salescreditmemo.dataMasterCustomer = ko.observableArray([])
salescreditmemo.dataDropDownCustomer = ko.observableArray([])
salescreditmemo.dataDropDownAccount = ko.observableArray([])
salescreditmemo.dataDropDownPO = ko.observableArray([])
salescreditmemo.dataMasterPurchaseOrder = ko.observableArray([])
salescreditmemo.dataMasterAccount = ko.observableArray([])
salescreditmemo.dataMastersalescreditmemo = ko.observableArray([])
salescreditmemo.dataMastersalescreditmemoOriginal = ko.observableArray([])
salescreditmemo.dataDropDownCustomerFilter = ko.observableArray([])
salescreditmemo.dataMasterInvoice = ko.observableArray([])
salescreditmemo.statusText = ko.observable()
salescreditmemo.textSearch = ko.observable()
salescreditmemo.filterindicator = ko.observable(false)
salescreditmemo.text = ko.observable()
salescreditmemo.TitelFilter = ko.observable(" Hide Filter")
salescreditmemo.DatePageBar = ko.observable()
salescreditmemo.textCustomerSearch = ko.observable()
salescreditmemo.showCreate = ko.observable(false)
salescreditmemo.backToList = ko.observable(false)
salescreditmemo.printPDFListView = ko.observable(false)
salescreditmemo.showEdit = ko.observable(false)
salescreditmemo.dataMasterInventory = ko.observableArray([])
salescreditmemo.dataDropDownInventory = ko.observableArray([])
salescreditmemo.warehouse = ko.observableArray([])
salescreditmemo.valueStorehouse = ko.observable()
salescreditmemo.dataLocation = ko.observableArray([])
salescreditmemo.dataDropDownSales = ko.observableArray([])
salescreditmemo.dataMasterSales = ko.observableArray([])

salescreditmemo.filterStatus = [{
    value: "DRAFT",
    text: "DRAFT"
}, {
    value: "POSTING",
    text: "POSTING"
}]

salescreditmemo.acccode = [{
    value: 5100,
    text: "SALES"
}, {
    value: 5200,
    text: "REVENUE"
}]

salescreditmemo.BoolVat = ko.observable()
salescreditmemo.sequenceNumber = ko.observable()
salescreditmemo.roman = [{
    1: "I"
}, {
    2: "II"
}, {
    3: "III"
}, {
    4: "IV"
}, {
    5: "V"
}, {
    6: "VI"
}, {
    7: "VII"
}, {
    8: "VIII"
}, {
    9: "Iresult"
}, {
    10: "result"
}, {
    11: "resultI"
}, {
    12: "resultII"
}]
salescreditmemo.newRecord = function () {
    var page = {
        ID: "",
        // AccountCode: 0,
        // AccountName: "",
        DocumentNo: "",
        DocumentNoInvoice: "",
        CustomerCode: "",
        CustomerName: "",
        DateCreated: "",
        DateStr: "",
        PoNumber: "",
        ListItem: [],
        Total: 0,
        VAT: 0,
        GrandTotalIDR: 0,
        GrandTotalUSD: 0,
        Rate: 1,
        Status: "",
        Description: "",
        Currency: "",
        StoreLocationId: 0,
        StoreLocationName: "",
        SalesCode: "",
        SalesName: "",
    }

    page.ListItem.push(salescreditmemo.listDetail({}))

    return page
}

salescreditmemo.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    salescreditmemo.DatePageBar(page)
}

salescreditmemo.listDetail = function (data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.CodeItem = data.Item == undefined ? '' : data.CodeItem
    dataTmp.Item = data.Item == undefined ? '' : data.Item
    dataTmp.Qty = data.Qty == undefined ? 0 : data.Qty
    dataTmp.PriceUSD = data.PriceUSD == undefined ? 0 : data.PriceUSD
    dataTmp.PriceIDR = data.PriceIDR == undefined ? 0 : data.PriceIDR
    dataTmp.AmountUSD = data.AmountUSD == undefined ? 0 : data.AmountUSD
    dataTmp.AmountIDR = data.AmountIDR == undefined ? 0 : data.AmountIDR

    var x = ko.mapping.fromJS(dataTmp)

    x.Qty.subscribe(function (e) {
        total = FormatCurrency(e) * FormatCurrency(x.PriceIDR())
        x.AmountIDR(total)

        //Menghitung Total Amount USD
        totalusd = FormatCurrency(e) * FormatCurrency(x.PriceUSD())
        x.AmountUSD(totalusd)

        //Menghitung Total 
        var alldataidr = ko.mapping.toJS(salescreditmemo.record.ListItem())
        var totalIDR = _.sumBy(alldataidr, function (v) {
            return v.AmountIDR
        })

        var alldata = ko.mapping.toJS(salescreditmemo.record.ListItem())
        var totalUSD = _.sumBy(alldata, function (v) {
            return v.AmountUSD
        })

        if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceIDR()) > 0) {
            salescreditmemo.record.Total(ChangeToRupiah(totalIDR))
        }

        if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceUSD()) > 0) {
            salescreditmemo.record.Total(ChangeToRupiah(totalUSD))
        }

    })

    x.PriceIDR.subscribe(function (e) {
        x.AmountIDR(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var alldataidr = ko.mapping.toJS(salescreditmemo.record.ListItem())
        var totalIDR = _.sumBy(alldataidr, function (v) {
            return v.AmountIDR
        })
        salescreditmemo.record.Total(ChangeToRupiah(totalIDR))

        salescreditmemo.checkdata()
    })

    x.PriceUSD.subscribe(function (e) {
        x.AmountUSD(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var alldatausd = ko.mapping.toJS(salescreditmemo.record.ListItem())
        var totalUSD = _.sumBy(alldatausd, function (v) {
            return v.AmountUSD
        })
        salescreditmemo.record.Total(ChangeToRupiah(totalUSD))
        salescreditmemo.checkdata()
    })


    return x
}


salescreditmemo.record = ko.mapping.fromJS(salescreditmemo.newRecord())

salescreditmemo.onChangePoNumber = function (value) {
    var data = _.filter(salescreditmemo.dataDropDownPO(), {
        'DocumentNumber': value
    })


    if (data.length != 0) {
        salescreditmemo.record.ListItem([])

        for (var i = 0; i < data[0].ListDetail.length; i++) {

            salescreditmemo.record.ListItem.push(salescreditmemo.listDetail({}))
            salescreditmemo.maskingMoney()

            salescreditmemo.record.ListItem()[i].Id(data[0].ListDetail[i].Id)
            salescreditmemo.record.ListItem()[i].Item(data[0].ListDetail[i].Item)
            salescreditmemo.record.ListItem()[i].Qty(data[0].ListDetail[i].Qty)
            salescreditmemo.record.ListItem()[i].PriceUSD(data[0].ListDetail[i].PriceUSD)
            salescreditmemo.record.ListItem()[i].PriceIDR(data[0].ListDetail[i].PriceIDR)
            salescreditmemo.record.ListItem()[i].AmountUSD(data[0].ListDetail[i].AmountUSD)
            salescreditmemo.record.ListItem()[i].AmountIDR(data[0].ListDetail[i].AmountIDR)
        }
    } else {
        salescreditmemo.record.ListItem([])
        salescreditmemo.record.ListItem.push(salescreditmemo.listDetail({}))
        salescreditmemo.record.Total(0)
    }


}

salescreditmemo.switchButton = function () {
    $('#my-checkbox').on('switchChange.bootstrapSwitch', function (event, state) {
        var data = ko.mapping.toJS(salescreditmemo.record)
        salescreditmemo.BoolVat(state)
        if (state) {
            var totVat = FormatCurrency(data.Total) / 10
            salescreditmemo.record.VAT(ChangeToRupiah(totVat))
        } else {
            salescreditmemo.record.VAT(0)
        }
    });
}

salescreditmemo.record.Total.subscribe(function (e) {
    var data = ko.mapping.toJS(salescreditmemo.record)
    if (salescreditmemo.BoolVat()) {
        var totVat = FormatCurrency(e) / 10
        salescreditmemo.record.VAT(ChangeToRupiah(totVat))
    } else {
        if (data.ListItem.length != 0) {
            if (FormatCurrency(data.ListItem[0].AmountUSD) > 0) {
                var GtUSD = FormatCurrency(data.Total)
                var GtIDR = FormatCurrency(data.Total) * FormatCurrency(data.Rate)
                salescreditmemo.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
                salescreditmemo.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
            } else {
                var GtIDR = FormatCurrency(data.Total)
                salescreditmemo.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
                if (data.Rate != 0) {
                    var GtUSD = GtIDR / FormatCurrency(data.Rate)
                    salescreditmemo.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
                } else {
                    salescreditmemo.record.GrandTotalUSD(0)
                }
            }
        } else {
            salescreditmemo.record.GrandTotalUSD(0)
            salescreditmemo.record.GrandTotalIDR(0)
        }
    }
})

salescreditmemo.record.VAT.subscribe(function (e) {
    var data = ko.mapping.toJS(salescreditmemo.record)
    if (data.ListItem.length != 0) {
        if (data.ListItem[0].AmountUSD != 0) {
            var GtUSD = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
            var GtIDR = GtUSD * FormatCurrency(data.Rate)
            // var GtIDR = (FormatCurrency(data.Total) / FormatCurrency(data.Rate)) + (FormatCurrency(data.VAT) / FormatCurrency(data.Rate))
            salescreditmemo.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
            salescreditmemo.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
        } else {
            var GtIDR = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
            salescreditmemo.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
            if (data.Rate != 0) {
                var GtUSD = GtIDR / FormatCurrency(data.Rate)
                salescreditmemo.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
            } else {
                salescreditmemo.record.GrandTotalUSD(0)
            }
        }
    } else {
        salescreditmemo.record.GrandTotalUSD(0)
        salescreditmemo.record.GrandTotalIDR(0)
    }
})

salescreditmemo.record.Rate.subscribe(function (e) {
    var data = ko.mapping.toJS(salescreditmemo.record)
    if (e <= 0) {
        return swal({
            title: "Error!",
            text: "Minimum value of rate is 1",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function () {
            salescreditmemo.record.Rate(1)
        });
    }
    if (data.ListItem[0].PriceUSD != 0) {
        var GtUSD = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
        var GtIDR = FormatCurrency(GtUSD) * FormatCurrency(e)
        salescreditmemo.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
        salescreditmemo.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
    } else {
        var GtIDR = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
        var GtUSD = FormatCurrency(GtIDR) / FormatCurrency(e)
        salescreditmemo.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
        salescreditmemo.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
    }
})

salescreditmemo.checkdata = function () {
    if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceUSD()) == 0 && FormatCurrency(salescreditmemo.record.ListItem()[0].PriceIDR()) == 0) {
        $(".priceidr").removeAttr("disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceUSD()) > 0) {
        $(".priceidr").attr("disabled", "disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceUSD()) == 0) {
        $(".priceusd").attr("disabled", "disabled");
        $(".priceidr").removeAttr("disabled");
    }
}

salescreditmemo.getDataSales = function() {
    model.Processing(true)
    ajaxPost('/master/getdatasales', {}, function(res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        salescreditmemo.dataMasterSales(res.Data)
        var DataSales = res.Data
        for (i in DataSales) {
            DataSales[i].Kode = DataSales[i].SalesID + ""
            DataSales[i].Name = DataSales[i].SalesID + " - " + DataSales[i].SalesName
        }
        salescreditmemo.dataDropDownSales(DataSales)

        model.Processing(false)
    })
}

salescreditmemo.getDataPurchaseOrder = function () {
    model.Processing(true)
    ajaxPost('/transaction/getpostingpurchaseorder', {}, function (res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        salescreditmemo.dataMasterPurchaseOrder(res.Data)
        salescreditmemo.dataDropDownPO(res.Data)
        salescreditmemo.processsPlusPlus()
        model.Processing(false)
    })
}

salescreditmemo.getDatasalescreditmemo = function (callback) {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {}
    if (salescreditmemo.filterindicator() == true) {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true,
            TextSearch: salescreditmemo.textSearch(),
            CustomerCode: salescreditmemo.textCustomerSearch(),
            Status: salescreditmemo.statusText(),
            LocationID: parseInt(userinfo.locationid())
        }
    } else {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: false,
            LocationID: parseInt(userinfo.locationid())
        }
    }
    if (salescreditmemo.valueStorehouse()!= "" && salescreditmemo.valueStorehouse()!= undefined && salescreditmemo.valueStorehouse()!= null) {
        param.LocationID = parseInt(salescreditmemo.valueStorehouse())
    }
    model.Processing(true)
    ajaxPost('/transaction/getdatasalescreditmemo', param, function (res) {
        if (res.IsError) {
            swal({
                title: "No Data Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            res.Data = []
            // return
        }
        model.Processing(false)
        salescreditmemo.dataMastersalescreditmemo(res.Data)
        salescreditmemo.dataMastersalescreditmemoOriginal(res.Data)
        salescreditmemo.processsPlusPlus()
        callback()
    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

salescreditmemo.getDataCustomer = function () {
    model.Processing(true)

    ajaxPost('/transaction/getcustomer', {}, function (res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        salescreditmemo.dataMasterCustomer(res.Data)
        salescreditmemo.dataDropDownCustomer(res.Data)
        var DataCustomer = res.Data
        for (i in DataCustomer) {
            DataCustomer[i].Kode = DataCustomer[i].Kode + ""
            DataCustomer[i].Name = DataCustomer[i].Kode + "-" + DataCustomer[i].Name
        }
        salescreditmemo.dataDropDownCustomerFilter(DataCustomer)
        salescreditmemo.processsPlusPlus()
        model.Processing(false)
    })
}

salescreditmemo.getDataAccount = function () {
    model.Processing(true)

    ajaxPost('/transaction/getaccount', {}, function (res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        salescreditmemo.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        salescreditmemo.dataDropDownAccount(DataAccount)
        salescreditmemo.processsPlusPlus()
        model.Processing(false)
    })
}

// salescreditmemo.onChangeCustomerCode = function (value) {
//     var result = _.filter(salescreditmemo.dataMasterCustomer(), {
//         'Kode': value
//     })[0].Name
//     salescreditmemo.record.CustomerName(result);
// }

salescreditmemo.onChangeCustomerCode = function (value) {
    if (value != "") {
        var result = _.filter(salescreditmemo.dataMasterCustomer(), {
            'Kode': value
        })[0].Name
        salescreditmemo.record.CustomerName(result);
    } else {
        salescreditmemo.record.CustomerName("");
    }

}

salescreditmemo.addNewItem = function () {
    salescreditmemo.record.ListItem.push(salescreditmemo.listDetail({}))
    salescreditmemo.maskingMoney()
    salescreditmemo.checkdata()
}

salescreditmemo.formCreated = function () {
    salescreditmemo.hideFilter()
    salescreditmemo.resetView()
    salescreditmemo.showCreate(true)
    salescreditmemo.showEdit(false)
    salescreditmemo.backToList(false)
    salescreditmemo.printPDFListView(false)
    salescreditmemo.disableView()
    salescreditmemo.record.ListItem([])
    salescreditmemo.record.ListItem.push(salescreditmemo.listDetail({}))
    salescreditmemo.maskingMoney()
    salescreditmemo.record.DateCreated(moment(new Date()).format('DD-MMM-YYYY'))
    $("#invnumber").removeAttr("disabled");
    $("#invnumber").val("")
    $("#Save").prop("disabled", true);
}

salescreditmemo.monthRoman = function (e) {
    var data = _.find(salescreditmemo.roman, e);
    return _.values(data)
}

salescreditmemo.reset = function () {
    salescreditmemo.record.ListItem([]);
    salescreditmemo.record.Total(0)
    salescreditmemo.resetView()
}

salescreditmemo.removerow = function () {
    salescreditmemo.checkdata()
    salescreditmemo.record.ListItem.remove(this)
    if (salescreditmemo.record.ListItem().length == 0) {
        salescreditmemo.record.ListItem.push(salescreditmemo.listDetail({}))
        salescreditmemo.record.Total(0)
    } else {
        //Menghitung Total 
        var alldataidr = ko.mapping.toJS(salescreditmemo.record.ListItem())
        var totalIDR = _.sumBy(alldataidr, function (v) {
            return v.AmountIDR
        })

        var alldata = ko.mapping.toJS(salescreditmemo.record.ListItem())
        var totalUSD = _.sumBy(alldata, function (v) {
            return v.AmountUSD
        })

        if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceIDR()) > 0) {
            salescreditmemo.record.Total(ChangeToRupiah(totalIDR))
        }

        if (FormatCurrency(salescreditmemo.record.ListItem()[0].PriceUSD()) > 0) {
            salescreditmemo.record.Total(ChangeToRupiah(totalUSD))
        }
    }
}

salescreditmemo.maskingMoney = function () {
    $('.currency').inputmask("numeric", {
        radiresultPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        rightAlign: false,
    });
}


salescreditmemo.renderGrid = function () {
    var mydata = salescreditmemo.dataMastersalescreditmemo()
    var currentUser = userinfo.usernameh()
    var currentRole = userinfo.rolenameh()
    for (var i = 0; i < mydata.length; i++) {
        var kondisi = false
        if (mydata[i].Status == 'DRAFT' && (mydata[i].User == currentUser || currentRole == 'administrator' || currentRole == 'supervisor')) {
            kondisi = true
        }
        mydata[i].Kondisi = kondisi
        // console.log(mydata[i].Kondisi)
    }
    // console.log(mydata)
    var columns = [
        {
            title: 'Action',
            width: 100,
            template: "#if (Kondisi) {# <button onclick='salescreditmemo.viewDraft(\"#: Id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #} else {# <button onclick='salescreditmemo.viewDraft(\"#: Id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>#}#"
        }, {
            field: 'DateCreated',
            title: 'Date Created',
            template: "#=moment(DateCreated).format('DD-MMM-YYYY')#",
            width: 160,
        },{
            field: 'StoreLocationName',
            title: 'Location',
            width: 160,
            template : function(dt){
                return dt.StoreLocationName + " (" + dt.StoreLocationId + ")"
            },
        },{
            field: 'DocumentNo',
            title: 'Doc No',
            width: 160,
        }, {
            field: 'CustomerName',
            title: 'Customer Name',
            width: 250,
        }, {
            field: 'GrandTotalIDR',
            title: 'Order Total',
            template: '#=ChangeToRupiah(FormatCurrency(GrandTotalIDR))#',
            width: 100,
            attributes: {
                "class": "rightAction",
            },
        }, {
            field: 'Status',
            title: 'Order Status',
            width: 100
        }]

    $('#gridListsalescreditmemo').kendoGrid({
        dataSource: {
            data: mydata,
            sort: {
                field: 'DateCreated',
                dir: 'desc',
            }
        },
        height: 500,
        width: 140,
        sortable: true,
        scrollable: {
            initialDirection: "desc"
        },
        columns: columns,
        excelExport: function (e) {
            ProActive.kendoExcelRender(e, "SalesCreditMemo", function (row, sheet) {
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "data") {
                        if (ci == 0) {
                            cell.value = moment(cell.value).toDate();
                            cell.format = "dd-MMM-yyyy";
                        }
                        if (ci == 3) {
                            cell.format = "#,##0.00_);(#,##0.00);0.00;";
                            cell.hAlign = "right";
                        }
                    }
                }
            });
        },
    })

}
salescreditmemo.exportExcel = function () {
    $("#gridListsalescreditmemo").getKendoGrid().saveAsExcel();
}

salescreditmemo.saveData = function () {
    var data = ko.mapping.toJS(salescreditmemo.record)

    data.DateCreated = $('#datesalescreditmemo').data('kendoDatePicker').value()

    for (var i = 0; i < data.ListItem.length; i++) {
        data.ListItem[i].Qty = parseInt(data.ListItem[i].Qty)
        data.ListItem[i].PriceIDR = FormatCurrency(data.ListItem[i].PriceIDR)
        data.ListItem[i].AmountIDR = FormatCurrency(data.ListItem[i].PriceIDR * data.ListItem[i].Qty)

        //in case user use IDR only
        if (FormatCurrency(data.ListItem[i].PriceUSD) == 0) {
            var USD = FormatCurrency(data.ListItem[i].PriceIDR) / FormatCurrency(data.Rate)
            data.ListItem[i].PriceUSD = USD
            data.ListItem[i].AmountUSD = USD * FormatCurrency(data.ListItem[i].Qty)
            data.Currency = "IDR"
        } else {
            data.Currency = "USD"
            data.ListItem[i].PriceUSD = FormatCurrency(data.ListItem[i].PriceUSD)
            data.ListItem[i].AmountUSD = FormatCurrency(data.ListItem[i].PriceUSD * data.ListItem[i].Qty)
        }
    }
    data.DateCreated = $('#datesalescreditmemo').data('kendoDatePicker').value()
    data.DateStr = moment(data.DateCreated).format('DD-MMM-YYYY')
    data.Status = "POSTING"
    data.Rate = FormatCurrency(data.Rate)
    // data.AccountCode = parseInt(salescreditmemo.record.AccountCode())

    // var c = _.result(_.find(salescreditmemo.acccode, function (obj) {
    //     return obj.value === parseInt(salescreditmemo.record.AccountCode());
    // }), 'text');
    // data.AccountName = c

    data.Total = FormatCurrency(data.Total)
    data.VAT = FormatCurrency(data.VAT)
    data.GrandTotalIDR = FormatCurrency(data.GrandTotalIDR)
    data.GrandTotalUSD = FormatCurrency(data.GrandTotalUSD)
    data.Description = data.Description

    if (data.Total == 0) {
        return swal({
            title: "Warning!",
            text: "None data for save",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.CustomerCode == "") {
        return swal({
            title: "Warning!",
            text: "No Customer Code has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.Rate == "") {
        return swal({
            title: "Warning!",
            text: "You haven't fill the Rate",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.Description == "") {
        return swal({
            title: "Warning!",
            text: "You haven't fill the Description",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else {
        for (var i = 0; i < data.ListItem.length; i++) {
            if (!(data.ListItem[i].AmountIDR || data.ListItem[i].AmountUSD)) {
                return swal({
                    title: "Warning!",
                    text: "Data in ListItem line hasn't completed yet",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            } else {
                var param = {
                    Data: data,
                    LastNumber: salescreditmemo.sequenceNumber(),
                }

                var url = "/transaction/insertsalescreditmemo"
                swal({
                    title: "Are you sure?",
                    text: "You will save this Sales Credit Memo",
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
                                    window.location.assign("/transaction/salescreditmemo")
                                    $("#Save").attr("disabled", "disabled")
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
        }
    }
}

salescreditmemo.delete = function () {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + ko.mapping.toJS(salescreditmemo.record).DocumentNo + "?",
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
            var url = "/transaction/deletesalescreditmemo";
            var param = {
                Id: ko.mapping.toJS(salescreditmemo.record).Id,
                DocumentNo: ko.mapping.toJS(salescreditmemo.record).DocumentNo,
            };
            ajaxPost(url, param, function (e) {
                if (e.Status == "OK") {
                    setTimeout(function () {
                        swal({
                            title: "Success!",
                            text: "Data has been deleted!",
                            type: "success",
                            confirmButtonColor: "#3da09a"
                        }, function () {
                            window.location.assign("/transaction/salescreditmemo")
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

salescreditmemo.PostingData = function () {
    var data = ko.mapping.toJS(salescreditmemo.record)

    data.DateCreated = $('#datesalescreditmemo').data('kendoDatePicker').value()

    data.Status = "POSTING"
    data.Rate = FormatCurrency(data.Rate)
    data.Total = FormatCurrency(data.Total)
    data.VAT = FormatCurrency(data.VAT)
    data.GrandTotalUSD = FormatCurrency(data.GrandTotalUSD)
    data.GrandTotalIDR = FormatCurrency(data.GrandTotalIDR)
    for (var i = 0; i < data.ListItem.length; i++) {
        data.ListItem[i].Qty = parseInt(data.ListItem[i].Qty)
        data.ListItem[i].PriceUSD = FormatCurrency(data.ListItem[i].PriceUSD)
        data.ListItem[i].PriceIDR = FormatCurrency(data.ListItem[i].PriceIDR)
        data.ListItem[i].AmountIDR = FormatCurrency(data.ListItem[i].AmountIDR)
        data.ListItem[i].AmountUSD = FormatCurrency(data.ListItem[i].AmountUSD)
    }

    if (data.Total == 0) {
        return swal({
            title: "Warning!",
            text: "None data for save",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.CustomerCode == "") {
        return swal({
            title: "Warning!",
            text: "No Customer Code has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.Rate == "") {
        return swal({
            title: "Warning!",
            text: "You havn't fill the Rate",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else {
        for (var i = 0; i < data.ListItem.length; i++) {
            if (!(data.ListItem[i].AmountIDR || data.ListItem[i].AmountUSD)) {
                return swal({
                    title: "Warning!",
                    text: "Data in ListItem line hasn't completed yet",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            } else {
                var param = {
                    Data: data,
                    LastNumber: salescreditmemo.sequenceNumber(),
                }

                var url = "/transaction/insertsalescreditmemo"
                swal({
                    title: "Are you sure?",
                    text: "You will Post this Sales Credit Memo",
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
                                    location.reload()
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
        }
    }
}

salescreditmemo.onChangeStatus = function (textFilter) {
    salescreditmemo.dataMastersalescreditmemo([])
    var allData = salescreditmemo.dataMastersalescreditmemoOriginal()
    if (textFilter != "" || textFilter != undefined) {

        var Data = _.filter(allData, function (o) {
            return o.Status.indexOf(textFilter) > -1
        });
        salescreditmemo.dataMastersalescreditmemo(Data)
    }
    salescreditmemo.renderGrid()
}

salescreditmemo.search = function (e) {
    salescreditmemo.textSearch(e)
    salescreditmemo.filterindicator(true)
    salescreditmemo.getDatasalescreditmemo(function () {
        salescreditmemo.renderGrid()
    })

}
salescreditmemo.printListToPdf = function (e) {
    model.Processing(true)
    var GrandTotalIDR = salescreditmemo.record.GrandTotalIDR()
    var total = FormatCurrency(GrandTotalIDR)
    var numberString = (total + "").split(".")[0];
    var comaString = (total + "").split(".")[1];
    var numWordStr = numToWords(numberString)
    var fixWord = numWordStr
    if (comaString != undefined) {
        var numwordComaStr = splitComa(comaString)
        fixWord = numWordStr + " point " + numwordComaStr
    }

    // console.log(numberString, comaString)
    // console.log(fixWord)
    var param = {
        Id: salescreditmemo.record.ID(),
        WordGrandtotal: fixWord,
    }
    ajaxPost("/transaction/exporttopdflistsalescreditmemo", param, function (e) {
        model.Processing(false)
        var taborWindow = window.open('/res/docs/salescreditmemo/' + e, '_blank');
        taborWindow.focus();
    })
}
function splitComa(str) {
    var word = ["zero", 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine']
    var array = str.split("")
    var len = array.length - 1
    var text = ""
    for (i in array) {
        var num = parseInt(array[i])
        if (i == len && num == 0) {
            text = text + ""
        } else {
            text = text + " " + word[num]
        }
    }
    return text
} 8
salescreditmemo.viewDraft = function (e) {
    salescreditmemo.hideFilter()
    var allData = salescreditmemo.dataMastersalescreditmemoOriginal()
    var data = _.find(allData, function (o) {
        return o.Id == e;
    });
    salescreditmemo.getDataInventory(data.StoreLocationId)
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    ko.mapping.fromJS(data, salescreditmemo.record)
    salescreditmemo.record.ID(data.Id)
    if (data.VAT > 0) {
        $('#my-checkbox').bootstrapSwitch('state', true);
    } else {
        $('#my-checkbox').bootstrapSwitch('state', false);
    }
    data.DateCreated = moment(data.DateCreated).format('DD-MMM-YYYY')
    data.Total = ChangeToRupiah(FormatCurrency(data.Total))
    data.VAT = ChangeToRupiah(FormatCurrency(data.VAT))
    data.GrandTotalIDR = ChangeToRupiah(FormatCurrency(data.GrandTotalIDR))
    data.GrandTotalUSD = ChangeToRupiah(FormatCurrency(data.GrandTotalUSD))
    data.Description = data.Description
    if (data.Currency == "USD") {
        for (var i = 0; i < data.ListItem.length; i++) {
            data.ListItem[i].PriceIDR = 0
            data.ListItem[i].AmountIDR = 0
        }
    } else {
        for (var i = 0; i < data.ListItem.length; i++) {
            data.ListItem[i].PriceUSD = 0
            data.ListItem[i].AmountUSD = 0
        }
    }
    // console.log(data)
    ko.mapping.fromJS(data, salescreditmemo.record)
    $("#listsalescreditmemo").removeClass("active")
    $("#invnumber").val(data.DocumentNoInvoice)
    $("#invnumber").attr("disabled", "disabled")
    $("#createsalescreditmemo").addClass("active")
    salescreditmemo.maskingMoney()
    salescreditmemo.text("View Sales Credit Memo")
    salescreditmemo.showCreate(false)
    salescreditmemo.showEdit(false)
    salescreditmemo.backToList(true)
    salescreditmemo.disableView()
    salescreditmemo.printPDFListView(true)
    setTimeout(function () {
        _.each(data.ListItem, function (v, i) {
            $("#itemcode_" + i).data("kendoDropDownList").value(v.CodeItem)
            drop = $("#itemcode_" + i).data("kendoDropDownList")
            drop.enable(false)
        })
    }, 100)
}

salescreditmemo.editDraft = function (e) {
    salescreditmemo.hideFilter()
    var allData = salescreditmemo.dataMastersalescreditmemoOriginal()
    var data = _.find(allData, function (o) {
        return o.Id == e;
    });
    data.DateCreated = moment(data.DateCreated).format('DD-MMM-YYYY HH:mm')
    $('#datesalescreditmemo').data('kendoDatePicker').value(new Date(data.DateCreated))
    data.Total = ChangeToRupiah(FormatCurrency(data.Total))
    data.VAT = ChangeToRupiah(FormatCurrency(data.VAT))
    data.GrandTotalIDR = ChangeToRupiah(FormatCurrency(data.GrandTotalIDR))
    data.GrandTotalUSD = ChangeToRupiah(FormatCurrency(data.GrandTotalUSD))
    if (data.Currency == "USD") {
        for (var i = 0; i < data.ListItem.length; i++) {
            data.ListItem[i].PriceIDR = 0
            data.ListItem[i].AmountIDR = 0
        }
    } else {
        for (var i = 0; i < data.ListItem.length; i++) {
            data.ListItem[i].PriceUSD = 0
            data.ListItem[i].AmountUSD = 0
        }
    }
    // data.Description = data.Description.substring(16)
    ko.mapping.fromJS(data, salescreditmemo.record)
    if (data.VAT > 0) {
        $('#my-checkbox').bootstrapSwitch('state', true);
    } else {
        $('#my-checkbox').bootstrapSwitch('state', false);
    }
    $("#listsalescreditmemo").removeClass("active")
    $("#createsalescreditmemo").addClass("active")
    salescreditmemo.record.ID(e)
    newArr = []
    _.each(salescreditmemo.record.ListItem(), function (v, i) {
        newArr.push(salescreditmemo.listDetail(ko.mapping.toJS(v)))
    })
    salescreditmemo.record.ListItem(newArr)
    salescreditmemo.maskingMoney()
    salescreditmemo.text("Edit Sales Credit Memo")
    salescreditmemo.showCreate(false)
    salescreditmemo.showEdit(true)
    salescreditmemo.backToList(true)
    salescreditmemo.printPDFListView(false)
    salescreditmemo.enableView()
    salescreditmemo.checkdata()
}

salescreditmemo.backList = function () {
    $('.nav-tabs a[href="#listsalescreditmemo"]').tab('show')
    $("#listsalescreditmemo").addClass("active");
    $("#createsalescreditmemo").removeClass("active");
}

salescreditmemo.disableView = function () {
    $(".formInput").attr("disabled", "disabled")
    $(".invhide").hide()
    // $(".btnDeleteSummary").attr("disabled", "disabled")
    // $("#addnewitem").attr("disabled", "disabled")
    $("#my-checkbox").bootstrapSwitch('disabled', true)
    $('#datesalescreditmemo').data('kendoDatePicker').enable(false);
    var dropDown2 = $("#customercode").data("kendoDropDownList");
    dropDown2.enable(false);
    var dropDown3 = $("#salescode").data("kendoDropDownList");
    dropDown3.enable(false);
    // var dropDown3 = $("#accountnumber").data("kendoDropDownList");
    // dropDown3.enable(false);
}

salescreditmemo.enableView = function () {
    $(".formInput").removeAttr("disabled")
    $(".invhide").show()
    // $(".btnDeleteSummary").removeAttr("disabled")
    // $("#addnewitem").removeAttr("disabled")
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#datesalescreditmemo').data('kendoDatePicker').enable(condition);
    var dropDown2 = $("#customercode").data("kendoDropDownList");
    dropDown2.enable(true);
    var dropDown3 = $("#salescode").data("kendoDropDownList");
    dropDown3.enable(true);
    // var dropDown3 = $("#accountnumber").data("kendoDropDownList");
    // dropDown3.enable(true);
    $("#my-checkbox").bootstrapSwitch({
        disabled: false
    });
}
salescreditmemo.loadUserLocation = function (callback) {
    ajaxPost('/master/getuserlocation', {}, function (res) {
        if (callback) callback.apply(this, [res]);
    });
}

salescreditmemo.resetView = function () {
    ko.mapping.fromJS(salescreditmemo.newRecord(), salescreditmemo.record)
    salescreditmemo.loadUserLocation(function (loc) {
        if (loc.IsError === false) {
            salescreditmemo.record.StoreLocationId(loc.Data.LocationID);
            salescreditmemo.record.StoreLocationName(loc.Data.LocationName);
            salescreditmemo.getDataInventory(loc.Data.LocationID);
        };
    });
    salescreditmemo.record.Rate(1)
    salescreditmemo.record.DateCreated(moment(new Date()).format('DD-MMM-YYYY'))
    $(".formInput").val("")
    $("#customername").val("")
    $(".Amount").val("")

    salescreditmemo.text("Create Sales Credit Memo")
    $('#customercode').data('kendoDropDownList').value(-1);
    // $('#accountnumber').data('kendoDropDownList').value(-1);
    salescreditmemo.enableView()
    $("#my-checkbox").bootstrapSwitch('state', false);
}

salescreditmemo.hideFilter = function () {
    var panelFilter = $('.panel-filter');
    var panelContent = $('.panel-content');
    panelFilter.hide();
    panelContent.attr('class', 'col-md-12 col-sm-12 ez panel-content');
    $('.breakdown-filter').removeAttr('style');
    salescreditmemo.TitelFilter(" Show Filter");
}

salescreditmemo.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
salescreditmemo.fromPOINVSummary = function () {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    if (num != null) {
        var allData = salescreditmemo.dataMastersalescreditmemoOriginal()
        var data = _.find(allData, function (o) {
            return o.DocumentNo == num;
        });
        // console.log(num, data)
        if (data != undefined) {
            salescreditmemo.viewDraft(data.Id)
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning"
            }, function () {
                window.location.assign("/transaction/salescreditmemo")
            });
        }
    }
}
salescreditmemo.processsPlusPlus = function () {
    processs += 1
    // console.log(processs)
    if (processs >= 4) {
        setTimeout(function () {
            salescreditmemo.fromPOINVSummary()
        }, 300)
        // setTimeut(purchasepayment.createPPFromPOInvSummary(),10000) 
    }
}
salescreditmemo.getDataInvoice = function (textData) {
    model.Processing(true)
    var param = {}
    param = {
        TextSearch: textData,
    }
    salescreditmemo.record.SalesCode("")
    salescreditmemo.record.SalesName("")
    salescreditmemo.record.CustomerCode("")
    salescreditmemo.record.CustomerName("")
    salescreditmemo.record.StoreLocationId("")
    salescreditmemo.record.StoreLocationName("")
    salescreditmemo.record.ListItem([])
    salescreditmemo.record.Total(0)
    salescreditmemo.record.VAT(0)
    salescreditmemo.record.GrandTotalIDR(0)
    salescreditmemo.record.Description("")
    salescreditmemo.disableView()
    $("#Save").prop("disabled", true);
    ajaxPost('/transaction/getdatainvoicesearch', param, function (res) {
        if (res.IsError) {
            swal({
                title: "Search Not Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            model.Processing(false)
            return
        }

        model.Processing(false)
        salescreditmemo.disableView()
        salescreditmemo.dataMasterInvoice(res.Data)
        var x = salescreditmemo.dataMasterInvoice()[0]
        salescreditmemo.getDataInventory(x.StoreLocationId)
        var arrCustName = x.CustomerName.split('-')
        //console.log(arrCustName[1])
        salescreditmemo.record.SalesCode(x.SalesCode)
        salescreditmemo.record.SalesName(x.SalesName)
        salescreditmemo.record.CustomerCode(x.CustomerCode)
        salescreditmemo.record.CustomerName(arrCustName[1])
        salescreditmemo.record.DocumentNoInvoice(x.DocumentNo)
        salescreditmemo.record.StoreLocationId(x.StoreLocationId)
        salescreditmemo.record.StoreLocationName(x.StoreLocationName)
        salescreditmemo.record.ListItem([])
        salescreditmemo.record.Description(x.DocumentNo + " - ");
        for (var i = 0; i < x.ListItem.length; i++) {
            salescreditmemo.record.ListItem.push(salescreditmemo.listDetail({}))
            salescreditmemo.record.ListItem()[i].CodeItem(x.ListItem[i].CodeItem)
            salescreditmemo.record.ListItem()[i].Item(x.ListItem[i].Item)
            salescreditmemo.record.ListItem()[i].Qty(x.ListItem[i].Qty)
            salescreditmemo.record.ListItem()[i].PriceIDR(x.ListItem[i].PriceIDR)
            salescreditmemo.record.ListItem()[i].AmountIDR(x.ListItem[i].AmountIDR)
        }
        salescreditmemo.disableView()
        $("#Description").prop("disabled", false).focus();

        salescreditmemo.maskingMoney();
        $("#invnumber").prop("disabled", true);
        $("#Save").prop("disabled", false);
        /*
        setTimeout(function () {
            _.each(x.ListItem, function (v, i) {
                $("#itemcode_" + i).data("kendoDropDownList").value(v.CodeItem)
                drop = $("#itemcode_" + i).data("kendoDropDownList")
                drop.enable(false)
            })
        }, 100)
        */
    })
    /*
    setTimeout(function () {
        $(".btnDeleteSummary").attr("disabled", "disabled")
        $("#addnewitem").attr("disabled", "disabled")
        $("#my-checkbox").bootstrapSwitch('disabled', true)
        var dropDown2 = $("#customercode").data("kendoDropDownList");
        dropDown2.enable(false);
        $(".formInput").attr("disabled", "disabled")
        $("#Description").removeAttr("disabled")
    }, 100)
    */
}

salescreditmemo.getDataInventory = function (val) {
    var param = {
        LocationId: val
    }
    model.Processing(true)

    ajaxPost('/master/getdatainventory', param, function (res) {


        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        salescreditmemo.dataMasterInventory(res.Data)
        var DataInventory = res.Data
        for (i in DataInventory) {
            DataInventory[i].Kode = DataInventory[i].INVID + ""
            DataInventory[i].Name = DataInventory[i].INVID + " - " + DataInventory[i].INVDesc
        }
        salescreditmemo.dataDropDownInventory(DataInventory)

        model.Processing(false)
    })
}

salescreditmemo.onChangeCodeItem = function (value, index) {
    findaccount = _.find(salescreditmemo.dataMasterInventory(), {
        INVID: value
    })
    //console.log(findaccount.INVDesc)
    salescreditmemo.record.ListItem()[index].Item(findaccount.INVDesc);
}

salescreditmemo.init = function () {
    var now = new Date()

    salescreditmemo.date(moment(now).format("DD MMMM YYYY"))
    // salescreditmemo.setDate()
    // salescreditmemo.initDate()
    salescreditmemo.getDataSales()
    salescreditmemo.getMasterLocation()
    salescreditmemo.getDataCustomer()
    salescreditmemo.getDataPurchaseOrder()
    salescreditmemo.getDataAccount()
    salescreditmemo.getDatasalescreditmemo(function () {
        salescreditmemo.renderGrid()
    })
    salescreditmemo.text("Create Sales Credit Memo")
    salescreditmemo.maskingMoney()
    salescreditmemo.switchButton()
    salescreditmemo.getDateNow()
    salescreditmemo.initAutoComplete();
    // salescreditmemo.getDataInventory()

}
function numToWords(number) {

    //Validates the number input and makes it a string
    if (typeof number === 'string') {
        number = parseInt(number, 10);
    }
    if (typeof number === 'number' && isFinite(number)) {
        number = number.toString(10);
    } else {
        return 'This is not a valid number';
    }

    //Creates an array with the number's digits and
    //adds the necessary amount of 0 to make it fully 
    //divisible by 3
    var digits = number.split('');
    while (digits.length % 3 !== 0) {
        digits.unshift('0');
    }


    //Groups the digits in groups of three
    var digitsGroup = [];
    var numberOfGroups = digits.length / 3;
    for (var i = 0; i < numberOfGroups; i++) {
        digitsGroup[i] = digits.splice(0, 3);
    }
    // console.log(digitsGroup); //debug

    //Change the group's numerical values to text
    var digitsGroupLen = digitsGroup.length;
    var numTxt = [
        [null, 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine'], //hundreds
        [null, 'ten', 'twenty', 'thirty', 'forty', 'fifty', 'sixty', 'seventy', 'eighty', 'ninety'], //tens
        [null, 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine'] //ones
    ];
    var tenthsDifferent = ['ten', 'eleven', 'twelve', 'thirteen', 'fourteen', 'fifteen', 'sixteen', 'seventeen', 'eighteen', 'nineteen'];

    // j maps the groups in the digitsGroup
    // k maps the element's position in the group to the numTxt equivalent
    // k values: 0 = hundreds, 1 = tens, 2 = ones
    for (var j = 0; j < digitsGroupLen; j++) {
        for (var k = 0; k < 3; k++) {
            var currentValue = digitsGroup[j][k];
            digitsGroup[j][k] = numTxt[k][currentValue];
            if (k === 0 && currentValue !== '0') { // !==0 avoids creating a string "null hundred"
                digitsGroup[j][k] += ' hundred ';
            } else if (k === 1 && currentValue === '1') { //Changes the value in the tens place and erases the value in the ones place
                digitsGroup[j][k] = tenthsDifferent[digitsGroup[j][2]];
                digitsGroup[j][2] = 0; //Sets to null. Because it sets the next k to be evaluated, setting this to null doesn't work.
            }
        }
    }

    // console.log(digitsGroup); //debug

    //Adds '-' for gramar, cleans all null values, joins the group's elements into a string
    for (var l = 0; l < digitsGroupLen; l++) {
        if (digitsGroup[l][1] && digitsGroup[l][2]) {
            digitsGroup[l][1] += '-';
        }
        digitsGroup[l].filter(function (e) { return e !== null });
        digitsGroup[l] = digitsGroup[l].join('');
    }

    // console.log(digitsGroup); //debug

    //Adds thousand, millions, billion and etc to the respective string.
    var posfix = [null, 'thousand', 'million', 'billion', 'trillion', 'quadrillion', 'quintillion', 'sextillion'];
    if (digitsGroupLen > 1) {
        var posfixRange = posfix.splice(0, digitsGroupLen).reverse();
        for (var m = 0; m < digitsGroupLen - 1; m++) { //'-1' prevents adding a null posfix to the last group
            if (digitsGroup[m]) {
                digitsGroup[m] += ' ' + posfixRange[m];
            }
        }
    }

    // console.log(digitsGroup); //debug

    //Joins all the string into one and returns it
    return digitsGroup.join(' ');

} //End of numToWords function
$(function () {
    salescreditmemo.init()
    $("#invnumber").on('keypress', function (e) {
        if (e.which == 13) {
            salescreditmemo.getDataInvoice($("#invnumber").val())
        }
    });
    $("#textSearch").on("keyup blur change", function () {
        salescreditmemo.filterText();
    });
})

// salescreditmemo.initDate = function () {
//     var dtpStart = $('#dateStart').data('kendoDatePicker');
//     var dtpEnd = $('#dateEnd').data('kendoDatePicker');
//     dtpStart.value(moment().startOf('month').toDate());
//     dtpEnd.value(moment().startOf('day').toDate());
//     dtpStart.max(dtpEnd.value());
//     dtpEnd.min(dtpStart.value());

//     dtpStart.bind("change", function () {
//         dtpEnd.min(dtpStart.value());
//     });
//     dtpEnd.bind("change", function () {
//         dtpStart.max(dtpEnd.value());
//     });
// }

salescreditmemo.onChangeDateStart = function(val){
    if (val.getTime()>salescreditmemo.DateEnd().getTime()){
        salescreditmemo.DateEnd(val)
    }
}

salescreditmemo.getMasterLocation= function() {
    $.ajax({
        url: "/transferorder/getuserlocations",
        success: function(data) {
            salescreditmemo.warehouse.removeAll();
            $(data).each(function(ix, ele) {
                salescreditmemo.dataLocation.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
                salescreditmemo.warehouse.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
            });
            salescreditmemo.valueStorehouse(data[0]["LocationID"])
        }
    });
}

salescreditmemo.filterText = function (term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["DocumentNo", "CustomerName", "Description", "Status", "StoreLocationName"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator: "contains", value: $searchValue });
    }
    $("#gridListsalescreditmemo").data("kendoGrid").dataSource.query({ filter: filter });
}

salescreditmemo.retrievedAutoCompleteData = [];
salescreditmemo.initAutoComplete = function() {
    salescreditmemo.autoCompleteDataSource = new kendo.data.DataSource({
        serverFiltering: true, 
        schema:
        {
            data: "data",
            count: "count"
        },
        transport:
        {
            read: function(e) {
                var filterValue = e.data.filter.filters[0] ? e.data.filter.filters[0].value : false;
                if (!filterValue) {
                    e.error();
                    return;
                }
                var url = "/transaction/getautocinvnumscm";
                ajaxPost(url, e.data.filter, function (res) {
                    salescreditmemo.retrievedAutoCompleteData = res.data;
                    e.success(res);
                }, function () {
                    e.error();
                })
            }
        }
    });
    $("#invnumber").kendoAutoComplete({
        dataSource: salescreditmemo.autoCompleteDataSource,
        filter: "contains",
        clearButton: false,
        placeholder: "Enter Invoice Number...",
        minLength: 12,
        dataValueField : "DocumentNo",
        dataTextField : "DocumentNo",
    });
}