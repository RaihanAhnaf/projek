model.Processing(false)

var invoice = {}
var processs = 0;
invoice.date = ko.observable()
invoice.dataMasterCustomer = ko.observableArray([])
invoice.dataDropDownCustomer = ko.observableArray([])
invoice.dataDropDownAccount = ko.observableArray([])
invoice.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
invoice.DateEnd = ko.observable(new Date)
invoice.dataDropDownPO = ko.observableArray([])
invoice.dataMasterPurchaseOrder = ko.observableArray([])
invoice.dataMasterAccount = ko.observableArray([])
invoice.dataMasterInvoice = ko.observableArray([])
invoice.dataMasterInvoiceOriginal = ko.observableArray([])
invoice.dataDropDownCustomerFilter = ko.observableArray([])
invoice.statusText = ko.observable()
invoice.textSearch = ko.observable()
invoice.filterindicator = ko.observable(false)
invoice.text = ko.observable()
invoice.TitelFilter = ko.observable(" Hide Filter")
invoice.DatePageBar = ko.observable()
invoice.textCustomerSearch = ko.observable()
invoice.showCreate = ko.observable(false)
invoice.backToList = ko.observable(false)
invoice.printPDFListView = ko.observable(false)
invoice.showEdit = ko.observable(false)
invoice.dataMasterInventory = ko.observableArray([])
invoice.dataDropDownInventory = ko.observableArray([])
invoice.dataDropDownSales = ko.observableArray([])
invoice.dataMasterSales = ko.observableArray([])
invoice.dataMasterLocation = ko.observableArray([])
invoice.dataMasterCategory = ko.observable("")
invoice.names = ko.observable()
invoice.dataMasterInvoiceNonInventory = ko.observableArray([])
invoice.dataMasterInvoiceOriginalNonInventory = ko.observableArray([])
invoice.warehouse = ko.observableArray([])
invoice.valueStorehouse = ko.observable()
invoice.dataLocation = ko.observableArray([])
invoice.isCustomerOK = ko.observable(0)
invoice.filterAlert = ko.observable(false) 
invoice.filterStatus = [{
    value: "DRAFT",
    text: "DRAFT"
}, {
    value: "POSTING",
    text: "POSTING"
}]

invoice.acccode = [{
    value: 5100,
    text: "SALES"
}, {
    value: 5200,
    text: "REVENUE"
}]

invoice.BoolVat = ko.observable()
invoice.sequenceNumber = ko.observable()
invoice.roman = [{
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
invoice.newRecord = function() {
    var page = {
        ID: "",
        AccountCode: 5110,
        AccountName: "PENJUALAN JASA / BARANG KENA PAJAK",
        DocumentNo: "",
        CustomerCode: "",
        CustomerName: "",
        DateCreated: "",
        DateStr: "",
        PoNumber: "",
        ListItem: [],
        Total: 0,
        VAT: 0,
        Discount: 0,
        GrandTotalIDR: 0,
        GrandTotalUSD: 0,
        Rate: 1,
        Status: "",
        Description: "",
        Currency: "",
        SalesCode: "",
        SalesName: "",
        StoreLocationId: 0,
        Category: "",
        StoreLocationName: "",
        INVCMI: false
    }
    page.ListItem.push(invoice.listDetail({}))

    return page
}

invoice.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    invoice.DatePageBar(page)
}

invoice.listDetail = function(data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
    dataTmp.Item = data.Item == undefined ? '' : data.Item
    dataTmp.Qty = data.Qty == undefined ? 0 : data.Qty
    dataTmp.PriceUSD = data.PriceUSD == undefined ? 0 : data.PriceUSD
    dataTmp.PriceIDR = data.PriceIDR == undefined ? 0 : data.PriceIDR
    dataTmp.AmountUSD = data.AmountUSD == undefined ? 0 : data.AmountUSD
    dataTmp.AmountIDR = data.AmountIDR == undefined ? 0 : data.AmountIDR

    var x = ko.mapping.fromJS(dataTmp)

    x.Qty.subscribe(function(e) {
        total = FormatCurrency(e) * FormatCurrency(x.PriceIDR())
        x.AmountIDR(total)

        //Menghitung Total Amount USD
        totalusd = FormatCurrency(e) * FormatCurrency(x.PriceUSD())
        x.AmountUSD(totalusd)

        //Menghitung Total 
        var alldataidr = ko.mapping.toJS(invoice.record.ListItem())
        var totalIDR = _.sumBy(alldataidr, function(v) {
            return v.AmountIDR
        })

        var alldata = ko.mapping.toJS(invoice.record.ListItem())
        var totalUSD = _.sumBy(alldata, function(v) {
            return v.AmountUSD
        })

        if (FormatCurrency(invoice.record.ListItem()[0].PriceIDR()) > 0) {
            invoice.record.Total(ChangeToRupiah(totalIDR))
        }

        if (FormatCurrency(invoice.record.ListItem()[0].PriceUSD()) > 0) {
            invoice.record.Total(ChangeToRupiah(totalUSD))
        }

        var invData = invoice.dataDropDownInventory();
        var itmCode = x.CodeItem();
        $(invData).each(function(idx, ele) {
            if (itmCode == ele.INVID) {
                if (x.Qty() > ele.Saldo) {
                    swal({
                        title: "Warning!",
                        text: "Item stock is less than qty!",
                        type: "info",
                        confirmButtonColor: "#3da09a"
                    })
                    x.Qty(0);
                }
            }
        });
    })

    x.PriceIDR.subscribe(function(e) {
        x.AmountIDR(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var alldataidr = ko.mapping.toJS(invoice.record.ListItem())
        var totalIDR = _.sumBy(alldataidr, function(v) {
            return v.AmountIDR
        })
        invoice.record.Total(ChangeToRupiah(totalIDR))

        invoice.checkdata()
    })

    x.PriceUSD.subscribe(function(e) {
        x.AmountUSD(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var alldatausd = ko.mapping.toJS(invoice.record.ListItem())
        var totalUSD = _.sumBy(alldatausd, function(v) {
            return v.AmountUSD
        })
        invoice.record.Total(ChangeToRupiah(totalUSD))
        invoice.checkdata()
    })


    return x
}


invoice.record = ko.mapping.fromJS(invoice.newRecord())

invoice.onChangePoNumber = function(value) {
    var data = _.filter(invoice.dataDropDownPO(), {
        'DocumentNumber': value
    })


    if (data.length != 0) {
        invoice.record.ListItem([])

        for (var i = 0; i < data[0].ListDetail.length; i++) {

            invoice.record.ListItem.push(invoice.listDetail({}))
            invoice.maskingMoney()

            invoice.record.ListItem()[i].Id(data[0].ListDetail[i].Id)
            invoice.record.ListItem()[i].CodeItem(data[0].ListDetail[i].CodeItem)
            invoice.record.ListItem()[i].Item(data[0].ListDetail[i].Item)
            invoice.record.ListItem()[i].Qty(data[0].ListDetail[i].Qty)
            invoice.record.ListItem()[i].PriceUSD(data[0].ListDetail[i].PriceUSD)
            invoice.record.ListItem()[i].PriceIDR(data[0].ListDetail[i].PriceIDR)
            invoice.record.ListItem()[i].AmountUSD(data[0].ListDetail[i].AmountUSD)
            invoice.record.ListItem()[i].AmountIDR(data[0].ListDetail[i].AmountIDR)
        }
    } else {
        invoice.record.ListItem([])
        invoice.record.ListItem.push(invoice.listDetail({}))
        invoice.record.Total(0)
    }


}

invoice.switchButton = function() {
    $('#my-checkbox').on('switchChange.bootstrapSwitch', function(event, state) {
        var data = ko.mapping.toJS(invoice.record)
        invoice.BoolVat(state)
        if (state) {
            // var totVat = FormatCurrency(data.Total) / 10
            var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.Total)
            invoice.record.VAT(ChangeToRupiah(totVat))
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.Total)
            console.log(GT, totVat)
            var Totalall = GT + totVat
            invoice.record.GrandTotalIDR(ChangeToRupiah(Totalall))
        } else {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.Total)
            invoice.record.GrandTotalIDR(ChangeToRupiah(GT))
            invoice.record.VAT(0)
        }
    });
}

invoice.record.Total.subscribe(function(e) {
    var data = ko.mapping.toJS(invoice.record)
    if (invoice.BoolVat()) {
        var totVat = FormatCurrency(e) / 10
        invoice.record.VAT(ChangeToRupiah(totVat))
    } else {
        if (data.ListItem.length != 0) {
            if (FormatCurrency(data.ListItem[0].AmountUSD) > 0) {
                var GtUSD = FormatCurrency(data.Total)
                var GtIDR = FormatCurrency(data.Total) * FormatCurrency(data.Rate)
                invoice.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
                invoice.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
            } else {
                var GtIDR = FormatCurrency(data.Total)
                invoice.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
                if (data.Rate != 0) {
                    var GtUSD = GtIDR / FormatCurrency(data.Rate)
                    invoice.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
                } else {
                    invoice.record.GrandTotalUSD(0)
                }
            }
        } else {
            invoice.record.GrandTotalUSD(0)
            invoice.record.GrandTotalIDR(0)
        }
    }
})

invoice.record.VAT.subscribe(function(e) {
    var data = ko.mapping.toJS(invoice.record)
    if (data.ListItem.length != 0) {
        if (data.ListItem[0].AmountUSD != 0) {
            var GtUSD = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
            var GtIDR = GtUSD * FormatCurrency(data.Rate)
            // var GtIDR = (FormatCurrency(data.Total) / FormatCurrency(data.Rate)) + (FormatCurrency(data.VAT) / FormatCurrency(data.Rate))
            invoice.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
            invoice.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
        } else {
            if (invoice.BoolVat()) {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.Total)
                invoice.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.Total)
                var Totalall = GT + totVat
                invoice.record.GrandTotalIDR(ChangeToRupiah(Totalall))
                // var GtIDR = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
                // invoice.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
                if (data.Rate != 0) {
                    var GtUSD = GtIDR / FormatCurrency(data.Rate)
                    invoice.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
                } else {
                    invoice.record.GrandTotalUSD(0)
                }
            }else{
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.Total)
                invoice.record.GrandTotalIDR(ChangeToRupiah(GT))
            }
        }
    } else {
        invoice.record.GrandTotalUSD(0)
        invoice.record.GrandTotalIDR(0)
    }
})

invoice.record.Rate.subscribe(function(e) {
    var data = ko.mapping.toJS(invoice.record)
    if (e <= 0) {
        return swal({
            title: "Error!",
            text: "Minimum value of rate is 1",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function() {
            invoice.record.Rate(1)
        });
    }
    if (data.ListItem[0].PriceUSD != 0) {
        var GtUSD = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
        var GtIDR = FormatCurrency(GtUSD) * FormatCurrency(e)
        invoice.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
        invoice.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
    } else {
        var GtIDR = FormatCurrency(data.Total) + FormatCurrency(data.VAT)
        var GtUSD = FormatCurrency(GtIDR) / FormatCurrency(e)
        invoice.record.GrandTotalUSD(ChangeToRupiah(GtUSD))
        invoice.record.GrandTotalIDR(ChangeToRupiah(GtIDR))
    }
})
invoice.record.Discount.subscribe(function (e) {
    var data = ko.mapping.toJS(invoice.record)
    if (FormatCurrency(e) > 100 || FormatCurrency(e) < 0) {
        return swal({
            title: "Warning!",
            text: "your discount are irational",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function () {
            invoice.record.Discount(0)
        });
    }
    var GT = ((100 - FormatCurrency(e)) / 100) * FormatCurrency(data.Total)
    var totVat = 0
    if (invoice.BoolVat()) {
        totVat = ((100 - FormatCurrency(e)) / 1000) * FormatCurrency(data.Total)
        invoice.record.VAT(ChangeToRupiah(totVat))
    } else {
        invoice.record.VAT(ChangeToRupiah(0))
        totVat = 0
    }
    invoice.record.GrandTotalIDR(ChangeToRupiah(GT + totVat))
})
invoice.checkdata = function() {
    if (FormatCurrency(invoice.record.ListItem()[0].PriceUSD()) == 0 && FormatCurrency(invoice.record.ListItem()[0].PriceIDR()) == 0) {
        $(".priceidr").removeAttr("disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(invoice.record.ListItem()[0].PriceUSD()) > 0) {
        $(".priceidr").attr("disabled", "disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(invoice.record.ListItem()[0].PriceUSD()) == 0) {
        $(".priceusd").attr("disabled", "disabled");
        $(".priceidr").removeAttr("disabled");
    }
}

invoice.getDataPurchaseOrder = function() {
    model.Processing(true)
    ajaxPost('/transaction/getpostingpurchaseorder', {}, function(res) {
        if (res.IsError) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        invoice.dataMasterPurchaseOrder(res.Data)
        invoice.dataDropDownPO(res.Data)
        invoice.processsPlusPlus()
        model.Processing(false)
    })
}

invoice.getDataInvoice = function(callback) {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {}
    if (invoice.filterindicator() == true) {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true,
            TextSearch: invoice.textSearch(),
            CustomerCode: invoice.textCustomerSearch(),
            Status: invoice.statusText(),
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
    if (invoice.valueStorehouse()!= "" && invoice.valueStorehouse()!= undefined && invoice.valueStorehouse()!= null) {
        param.LocationID = parseInt(invoice.valueStorehouse())
    }
    model.Processing(true)
    ajaxPost('/transaction/getdatainvoice', param, function(res) {
        model.Processing(false)
        if (res.IsError && invoice.filterAlert()) {
            swal({
                title: "Search Not Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            res.Data = []
            // return
        }
        invoice.dataMasterInvoice(res.Data)
        invoice.dataMasterInvoiceOriginal(res.Data)
        invoice.processsPlusPlus()
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

invoice.getDataInvoiceNonInventory = function(callback) {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {}
    if (invoice.filterindicator() == true) {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true,
            TextSearch: invoice.textSearch(),
            CustomerCode: invoice.textCustomerSearch(),
            Status: invoice.statusText(),
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
    if (invoice.valueStorehouse()!= "" && invoice.valueStorehouse()!= undefined && invoice.valueStorehouse()!= null) {
        param.LocationID = parseInt(invoice.valueStorehouse())
    }
    // console.log(param);
    model.Processing(true)
    ajaxPost('/transaction/getdatainvoicenoninventory', param, function(res) {
        model.Processing(false)
        if (res.IsError && invoice.filterAlert()) {
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
        invoice.dataMasterInvoiceNonInventory(res.Data)
        invoice.dataMasterInvoiceOriginalNonInventory(res.Data)
        invoice.processsPlusPlus()
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
invoice.getDataCustomer = function() {
    model.Processing(true)

    ajaxPost('/transaction/getcustomer', {}, function(res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        invoice.dataMasterCustomer(res.Data)
        invoice.dataDropDownCustomer(res.Data)
        var DataCustomer = res.Data
        for (i in DataCustomer) {
            DataCustomer[i].Kode = DataCustomer[i].Kode + ""
            DataCustomer[i].Name = DataCustomer[i].Kode + "-" + DataCustomer[i].Name
        }
        invoice.dataDropDownCustomerFilter(DataCustomer)
        invoice.processsPlusPlus()
        model.Processing(false)
    })
}

invoice.getDataAccount = function() {
    model.Processing(true)

    ajaxPost('/transaction/getaccount', {}, function(res) {

        if (res.Total === 0) {
            //console.log(res)
            swal({
                title: "Info!",
                text: "Chart of account is empty, Please input chart of account!",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        invoice.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        invoice.dataDropDownAccount(DataAccount)
        invoice.processsPlusPlus()
        model.Processing(false)
    })
}

invoice.onChangeCustomerCode = function(value) {
    console.log("here")
    var result = _.filter(invoice.dataMasterCustomer(), {
        'Kode': value
    })
    console.log(result)
    invoice.record.CustomerName(result[0].Name);
}

invoice.getDataSales = function() {
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
        invoice.dataMasterSales(res.Data)
        var DataSales = res.Data
        for (i in DataSales) {
            DataSales[i].Kode = DataSales[i].SalesID + ""
            DataSales[i].Name = DataSales[i].SalesID + " - " + DataSales[i].SalesName
        }
        invoice.dataDropDownSales(DataSales)

        model.Processing(false)
    })
}

invoice.onChangeCustomerCode = function(value) {
    if (value != "") {
        var result = _.filter(invoice.dataMasterCustomer(), {
            'Kode': value
        })
        var param = {
            CustomerCode: value,
            DueDate     : moment().subtract(result[0].TrxCode,"days").format("YYYY-MM-DD"),
            DateNow   : moment().format("YYYY-MM-DD")
        };
        ajaxPost('/transaction/checkavailablecustomer', param, function(res) {
            var result = _.filter(invoice.dataMasterCustomer(), {
                'Kode': value
            })
            if (res.Data["BalanceIDR"] == 0) {
                invoice.record.CustomerName(result[0].Name);
            } else {
                if (res.Data["BalanceIDR"] > 0) {
                    swal({
                        title: "Unable to proceed customer!",
                        text: result[0].Name+" : \nThere is unpaid transaction out of due date",
                        type: "error",
                        confirmButtonColor: "#3da09a"
                    })
                } else {
                    swal({
                        title: "Unable to proceed this customer!",
                        text: result[0].Name+" : \nThe transaction of the customer has reach out of balance limit",
                        type: "error",
                        confirmButtonColor: "#3da09a"
                    })
                }
                setTimeout(() => {
                    $('#customercode').data('kendoDropDownList').value();
                    invoice.record.CustomerName(result[0].Name);
                }, 100);
                
            }
        })
    } else {
        invoice.record.CustomerName("");
    }
}

invoice.checkAvailableCustomer = function(customer_code) {
    var flag = 0;
    var result = _.filter(invoice.dataMasterCustomer(), {
        'Kode': customer_code
    })
    var param = {
        CustomerCode: customer_code,
        // DueDate     : moment().subtract(result[0].TrxCode,"days").format("yyyy-mm-dd"),
        DueDate     : moment().subtract(1,"days").format("YYYY-MM-DD"),
        DateNow   : moment().format("YYYY-MM-DD")
    };
    ajaxPost('/transaction/checkavailablecustomer', param, function(res) {
        if (res.Data["BalanceIDR"] == 0) {
            
        } else {
            
        }
    })
}

invoice.addNewItem = function() {
    invoice.record.ListItem.push(invoice.listDetail({}))
    invoice.maskingMoney()
    invoice.checkdata()
}

invoice.formCreated = function() {
    invoice.hideFilter()
    invoice.resetView()
    invoice.showCreate(true)
    invoice.showEdit(false)
    invoice.dropdownCategory()
    invoice.backToList(false)
    invoice.printPDFListView(false)
    invoice.enableView()
    var loc = parseInt(userinfo.locationid() % 1000)
    if (userinfo.locationid() != "1000" && loc != 0) {
        invoice.record.StoreLocationId(parseInt(userinfo.locationid()))
        invoice.onChangeLocation(userinfo.locationid())
    }

    invoice.showEditor()
    invoice.record.DateCreated(moment(new Date()).format('DD-MMM-YYYY'))
    invoice.record.ListItem([])
    invoice.record.ListItem.push(invoice.listDetail({}))
    invoice.maskingMoney()
    invoice.isViewOnly = false;
    if (invoice.names() == "Non Inventory") {
        $(".invhideNon").hide()
        $(".item").prop('disabled', false);

    } else {
        $(".invhideNon").show()
        $(".item").prop('disabled', true);
    }
}

invoice.monthRoman = function(e) {
    var data = _.find(invoice.roman, e);
    return _.values(data)
}

invoice.reset = function() {
    invoice.record.ListItem([]);
    invoice.record.Total(0)
    invoice.resetView()
}

invoice.removerow = function() {
    invoice.checkdata()
    invoice.record.ListItem.remove(this)
    if (invoice.record.ListItem().length == 0) {
        invoice.record.ListItem.push(invoice.listDetail({}))
        invoice.record.Total(0)
    } else {
        //Menghitung Total 
        var alldataidr = ko.mapping.toJS(invoice.record.ListItem())
        var totalIDR = _.sumBy(alldataidr, function(v) {
            return v.AmountIDR
        })

        var alldata = ko.mapping.toJS(invoice.record.ListItem())
        var totalUSD = _.sumBy(alldata, function(v) {
            return v.AmountUSD
        })

        if (FormatCurrency(invoice.record.ListItem()[0].PriceIDR()) > 0) {
            invoice.record.Total(ChangeToRupiah(totalIDR))
        }

        if (FormatCurrency(invoice.record.ListItem()[0].PriceUSD()) > 0) {
            invoice.record.Total(ChangeToRupiah(totalUSD))
        }
        var data = ko.mapping.toJS(purchaseorder.record)
        //Menghitung VAT           
        var nomvat = ((100 - data.Discount) / 1000) * FormatCurrency(data.Total)
        invoice.record.VAT(ChangeToRupiah(nomvat))
        // Menghitung GrandTotal
        var totalnumber = FormatCurrency(data.Total)
        var nomgrandtotal = (((100 - data.Discount) / 100) * totalnumber) + nomvat
        invoice.record.GrandTotalIDR(ChangeToRupiah(nomgrandtotal))
    }
}

invoice.maskingMoney = function() {
    $('.currency').inputmask("numeric", {
        radiresultPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        rightAlign: false,
    });
}


invoice.renderGrid = function() {
    var mydata = invoice.dataMasterInvoice()
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
    var columns = [{
        title: 'Action',
        width: 100,
        template: "#if (Kondisi) {# <button onclick='invoice.viewDraft(\"#: Id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button onclick='invoice.editDraft(\"#: Id #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button> #} else {# <button onclick='invoice.viewDraft(\"#: Id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button onclick='' class='btn btn-sm btn-success btn-flat pointer-disable'><i class='fa fa-pencil'></i></button> #}#"
    }, {
        field: 'DateCreated',
        title: 'Date Created',
        template: "#=moment(DateCreated).format('DD-MMM-YYYY')#",
        width: 160,
    }, {
        field: 'StoreLocationName',
        title: 'Location',
        width: 160,
        template : function(dt){
            return dt.StoreLocationName + " (" + dt.StoreLocationId + ")"
        },
    }, {
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

    $('#gridListInvoice').kendoGrid({
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
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "Invoice", function(row, sheet) {
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

invoice.renderGridNonInventory = function() {
    var mydata = invoice.dataMasterInvoiceNonInventory()
    var currentUser = userinfo.usernameh()
    var currentRole = userinfo.rolenameh()
    if (mydata != undefined) {
        for (var i = 0; i < mydata.length; i++) {
            var kondisi = false
            if (mydata[i].Status == 'DRAFT' && (mydata[i].User == currentUser || currentRole == 'administrator' || currentRole == 'supervisor')) {
                kondisi = true
            }
            mydata[i].Kondisi = kondisi
            // console.log(mydata[i].Kondisi)
        }
    }
    // console.log(mydata)
    var columns = [{
        title: 'Action',
        width: 100,
        template: "#if (Kondisi) {# <button onclick='invoice.viewDraft(\"#: Id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button onclick='invoice.editDraft(\"#: Id #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button> #} else {# <button onclick='invoice.viewDraft(\"#: Id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button onclick='' class='btn btn-sm btn-success btn-flat pointer-disable'><i class='fa fa-pencil'></i></button> #}#"
    }, {
        field: 'DateCreated',
        title: 'Date Created',
        template: "#=moment(DateCreated).format('DD-MMM-YYYY')#",
        width: 160,
    },
    // {
    //     field: 'StoreLocationName',
    //     title: 'Location',
    //     template : function(dt){
    //         return dt.StoreLocationName + " (" + dt.StoreLocationId + ")"
    //     },
    //     width: 160,
    // }, 
    {
        field: 'DocumentNo',
        title: 'Doc No',
        width: 160,
        template : function(dt){
            var docs = dt.DocumentNo.split("/");
            if (docs.length>3){
                return docs[0]+"/"+docs[2]+"/"+docs[3]
            }
            return ""
        },
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
    $('#gridListInvoice').html("")
    $('#gridListInvoice').kendoGrid({
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
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "Invoice Non Inventory", function(row, sheet) {
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
invoice.exportExcel = function() {
    $("#gridListInvoice").getKendoGrid().saveAsExcel();
}

invoice.saveData = function() {
    var data = ko.mapping.toJS(invoice.record)

    data.DateCreated = $('#dateinvoice').data('kendoDatePicker').value()

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

        var invData = invoice.dataDropDownInventory();
        var itmCode = data.ListItem[i].CodeItem;
        var stockInvalid = false;
        $(invData).each(function(idx, ele) {
            if (itmCode == ele.INVID) {
                if (data.ListItem[i].Qty > ele.Saldo) {
                    stockInvalid = true;
                }
            }
        });
        if (stockInvalid)
            return swal({
                title: "Warning!",
                text: "Item #" + (i + 1) + ": stock is less than qty!",
                type: "info",
                confirmButtonColor: "#3da09a"
            });
    }
    data.DateCreated = $('#dateinvoice').data('kendoDatePicker').value()
    data.DateStr = moment(data.DateCreated).format('DD-MMM-YYYY')
    data.Status = "DRAFT"
    data.Rate = FormatCurrency(data.Rate)
    data.AccountCode = parseInt(invoice.record.AccountCode())
    data.StoreLocationId = parseInt(invoice.record.StoreLocationId())
    // var c = _.result(_.find(invoice.acccode, function(obj) {
    //     return obj.value === parseInt(invoice.record.AccountCode());
    // }), 'text');
    data.AccountName = data.AccountName

    tempGrandTotalIDR = data.GrandTotalIDR
    data.Total = FormatCurrency(data.Total)
    data.Discount = FormatCurrency(data.Discount)
    data.VAT = FormatCurrency(data.VAT)
    data.GrandTotalIDR = FormatCurrency(data.GrandTotalIDR)
    data.GrandTotalUSD = FormatCurrency(data.GrandTotalUSD)
    data.Description = data.Description
    data.SalesCode = data.SalesCode
    data.SalesName = data.SalesName
    if (data.Discount == 0) {
        data.GrandTotalIDR = data.Total + data.VAT
    }
    data.INVCMI = $('#checkvat-inv-cmi').bootstrapSwitch('state');
    if (data.INVCMI) {
        data.Status = "POSTING"
    }

    var checklimit = []
    if (data.CustomerCode != "") {
        checklimit = _.filter(invoice.dataMasterCustomer(), {
            'Kode': data.CustomerCode
        })
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
    } else {
        for (var i = 0; i < data.ListItem.length; i++) {
            if (!(data.ListItem[i].AmountIDR || data.ListItem[i].AmountUSD)) {
                return swal({
                    title: "Warning!",
                    text: "Data in ListItem line hasn't completed yet",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            }
        }
    }

    if(checklimit.length >0){
        var p = {
            CustomerCode : data.CustomerCode
        }
        ajaxPost('/master/getcustomerbalance', p, function (res) {
            if (!res.IsError) {
                checklimit[0].CurrentBalance = res.Data.BalanceIDR
            } else{
                checklimit[0].CurrentBalance = 0
            }
            if (( (FormatCurrency(tempGrandTotalIDR)+parseInt(checklimit[0].CurrentBalance)) > checklimit[0].Limit)) {
                return swal({
                    title: "Customer reach out balance limit!",
                    text: "Total Amount : Rp. "+ChangeToRupiah(data.GrandTotalIDR)+"\nCurrent Balance : Rp. "+ChangeToRupiah(parseInt(checklimit[0].CurrentBalance)) +"\nLimit : Rp. "+ChangeToRupiah(checklimit[0].Limit),
                    type: "warning",
                    confirmButtonColor: "#3da09a"
                })   
            }
            
        })
    }
    var param = {
        Data: data,
        LastNumber: invoice.sequenceNumber(),
    }

    var url = "/transaction/insertinvoice"
    swal({
        title: "Are you sure?",
        text: "You will save this Invoice",
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
                setTimeout(function() {
                    swal({
                        title: "Success!",
                        text: "Data has been saved!",
                        type: "success",
                        confirmButtonColor: "#3da09a"
                    }, function() {
                        invoice.resetView()
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
invoice.saveDataNonInventory = function() {
    var data = ko.mapping.toJS(invoice.record)

    data.DateCreated = $('#dateinvoice').data('kendoDatePicker').value()

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
    data.DateCreated = $('#dateinvoice').data('kendoDatePicker').value()
    data.DateStr = moment(data.DateCreated).format('DD-MMM-YYYY')
    data.Status = "DRAFT"
    data.Rate = FormatCurrency(data.Rate)
    data.AccountCode = parseInt(invoice.record.AccountCode())
    data.StoreLocationId = parseInt(invoice.record.StoreLocationId())
    data.AccountName = data.AccountName
    data.Category = $("#categoryDropdown").data("kendoDropDownList").value()//invoice.dataMasterCategory()
    tempGrandTotalIDR = data.GrandTotalIDR
    data.Total = FormatCurrency(data.Total)
    data.Discount = FormatCurrency(data.Discount)
    data.VAT = FormatCurrency(data.VAT)
    data.GrandTotalIDR = FormatCurrency(data.GrandTotalIDR)
    data.GrandTotalUSD = FormatCurrency(data.GrandTotalUSD)
    data.Description = data.Description
    data.SalesCode = data.SalesCode
    data.SalesName = data.SalesName
    if (data.Discount == 0) {
        data.GrandTotalIDR = data.Total + data.VAT
    }
    var checklimit = []
    if (data.CustomerCode != "") {
        checklimit = _.filter(invoice.dataMasterCustomer(), {
            'Kode': data.CustomerCode
        })
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
    } else {
        for (var i = 0; i < data.ListItem.length; i++) {
            if (!(data.ListItem[i].AmountIDR || data.ListItem[i].AmountUSD)) {
                return swal({
                    title: "Warning!",
                    text: "Data in ListItem line hasn't completed yet",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            }
        }
    }

    if(checklimit.length >0){
        var p = {
            CustomerCode : data.CustomerCode
        }
        ajaxPost('/master/getcustomerbalance', p, function (res) {
            if (!res.IsError) {
                checklimit[0].CurrentBalance = res.Data.BalanceIDR
            } else{
                checklimit[0].CurrentBalance = 0
            }
            if (( (FormatCurrency(tempGrandTotalIDR)+parseInt(checklimit[0].CurrentBalance)) > checklimit[0].Limit)) {
                return swal({
                    title: "Customer reach out balance limit!",
                    text: "Total Amount : Rp. "+ChangeToRupiah(data.GrandTotalIDR)+"\nCurrent Balance : Rp. "+ChangeToRupiah(parseInt(checklimit[0].CurrentBalance)) +"\nLimit : Rp. "+ChangeToRupiah(checklimit[0].Limit),
                    type: "warning",
                    confirmButtonColor: "#3da09a"
                })   
            }
            
        })
    }
    var param = {
        Data: data,
        LastNumber: invoice.sequenceNumber(),
    }

    var url = "/transaction/insertinvoicenoninventory"
    swal({
        title: "Are you sure?",
        text: "You will save this Invoice",
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
                setTimeout(function() {
                    swal({
                        title: "Success!",
                        text: "Data has been saved!",
                        type: "success",
                        confirmButtonColor: "#3da09a"
                    }, function() {
                        invoice.resetView()
                        invoice.formCreated()
                    });
                }, 100)
                console.log(e)
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

invoice.delete = function() {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + ko.mapping.toJS(invoice.record).DocumentNo + "?",
        text: "Your will not be able to recover this data",
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
            if (invoice.names() == "Non Inventory") {
                url = "/transaction/deleteinvoicenon";
            } else{
                url = "/transaction/deleteinvoice";
            }
            var param = {
                Id: ko.mapping.toJS(invoice.record).Id,
                DocumentNo: ko.mapping.toJS(invoice.record).DocumentNo,
            };
            ajaxPost(url, param, function(e) {
                if (e.Status == "OK") {
                    setTimeout(function() {
                        swal({
                            title: "Success!",
                            text: "Data has been deleted!",
                            type: "success",
                            confirmButtonColor: "#3da09a"
                        }, function() {
                            window.location.assign("/transaction/invoice")
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

invoice.PostingData = function() {
    var data = ko.mapping.toJS(invoice.record)

    data.DateCreated = $('#dateinvoice').data('kendoDatePicker').value()

    data.Status = "POSTING"
    data.Rate = FormatCurrency(data.Rate)
    data.Total = FormatCurrency(data.Total)
    data.VAT = FormatCurrency(data.VAT)
    tempGrandTotalIDR = data.GrandTotalIDR
    data.GrandTotalUSD = FormatCurrency(data.GrandTotalUSD)
    data.GrandTotalIDR = FormatCurrency(data.GrandTotalIDR)
    for (var i = 0; i < data.ListItem.length; i++) {
        data.ListItem[i].Qty = parseInt(data.ListItem[i].Qty)
        data.ListItem[i].PriceUSD = FormatCurrency(data.ListItem[i].PriceUSD)
        data.ListItem[i].PriceIDR = FormatCurrency(data.ListItem[i].PriceIDR)
        data.ListItem[i].AmountIDR = FormatCurrency(data.ListItem[i].AmountIDR)
        data.ListItem[i].AmountUSD = FormatCurrency(data.ListItem[i].AmountUSD)

        var invData = invoice.dataDropDownInventory();
        var itmCode = data.ListItem[i].CodeItem;
        var stockInvalid = false;
        $(invData).each(function(idx, ele) {
            if (itmCode == ele.INVID) {
                if (data.ListItem[i].Qty > ele.Saldo) {
                    stockInvalid = true;
                }
            }
        });
        if (stockInvalid)
            return swal({
                title: "Warning!",
                text: "Item #" + (i + 1) + ": stock is less than qty!",
                type: "info",
                confirmButtonColor: "#3da09a"
            });
    }

    var checklimit = []
    if (data.CustomerCode != "") {
        checklimit = _.filter(invoice.dataMasterCustomer(), {
            'Kode': data.CustomerCode
        })
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
            }
        }
    }

    if(checklimit.length >0){
        var p = {
            CustomerCode : data.CustomerCode
        }
        ajaxPost('/master/getcustomerbalance', p, function (res) {
            if (!res.IsError) {
                checklimit[0].CurrentBalance = res.Data.BalanceIDR
            } else{
                checklimit[0].CurrentBalance = 0
            }
            if (( (FormatCurrency(tempGrandTotalIDR)+parseInt(checklimit[0].CurrentBalance)) > checklimit[0].Limit)) {
                return swal({
                    title: "Customer reach out balance limit!",
                    text: "Total Amount : Rp. "+ChangeToRupiah(data.GrandTotalIDR)+"\nCurrent Balance : Rp. "+ChangeToRupiah(parseInt(checklimit[0].CurrentBalance)) +"\nLimit : Rp. "+ChangeToRupiah(checklimit[0].Limit),
                    type: "warning",
                    confirmButtonColor: "#3da09a"
                })   
            }
            
        })
    }

    var param = {
        Data: data,
        LastNumber: invoice.sequenceNumber(),
    }

    var url = "/transaction/insertinvoice"
    swal({
        title: "Are you sure?",
        text: "You will Post this Invoice",
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
                setTimeout(function() {
                    swal({
                        title: "Success!",
                        text: "Data has been saved!",
                        type: "success",
                        confirmButtonColor: "#3da09a"
                    }, function() {
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

invoice.PostingDataNonInventory = function() {
    var data = ko.mapping.toJS(invoice.record)

    data.DateCreated = $('#dateinvoice').data('kendoDatePicker').value()

    data.Status = "POSTING"
    data.Rate = FormatCurrency(data.Rate)
    data.Total = FormatCurrency(data.Total)
    data.VAT = FormatCurrency(data.VAT)
    tempGrandTotalIDR = data.GrandTotalIDR
    data.GrandTotalUSD = FormatCurrency(data.GrandTotalUSD)
    data.GrandTotalIDR = FormatCurrency(data.GrandTotalIDR)
    for (var i = 0; i < data.ListItem.length; i++) {
        data.ListItem[i].Qty = parseInt(data.ListItem[i].Qty)
        data.ListItem[i].PriceUSD = FormatCurrency(data.ListItem[i].PriceUSD)
        data.ListItem[i].PriceIDR = FormatCurrency(data.ListItem[i].PriceIDR)
        data.ListItem[i].AmountIDR = FormatCurrency(data.ListItem[i].AmountIDR)
        data.ListItem[i].AmountUSD = FormatCurrency(data.ListItem[i].AmountUSD)
    }

    var checklimit = []
    if (data.CustomerCode != "") {
        checklimit = _.filter(invoice.dataMasterCustomer(), {
            'Kode': data.CustomerCode
        })
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
            text: "You haven't fill the Rate",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.SalesCode == "") {
        return swal({
            title: "Warning!",
            text: "No Sales Code has selected",
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
            }
        }
    }

    if(checklimit.length >0){
        var p = {
            CustomerCode : data.CustomerCode
        }
        ajaxPost('/master/getcustomerbalance', p, function (res) {
            if (!res.IsError) {
                checklimit[0].CurrentBalance = res.Data.BalanceIDR
            } else{
                checklimit[0].CurrentBalance = 0
            }
            if (( (FormatCurrency(tempGrandTotalIDR)+parseInt(checklimit[0].CurrentBalance)) > checklimit[0].Limit)) {
                return swal({
                    title: "Customer reach out balance limit!",
                    text: "Total Amount : Rp. "+ChangeToRupiah(data.GrandTotalIDR)+"\nCurrent Balance : Rp. "+ChangeToRupiah(parseInt(checklimit[0].CurrentBalance)) +"\nLimit : Rp. "+ChangeToRupiah(checklimit[0].Limit),
                    type: "warning",
                    confirmButtonColor: "#3da09a"
                })   
            }
            
        })
    }

    var param = {
        Data: data,
        LastNumber: invoice.sequenceNumber(),
    }

    var url = "/transaction/insertinvoicenoninventory"
    swal({
        title: "Are you sure?",
        text: "You will Post this Invoice",
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
                setTimeout(function() {
                    swal({
                        title: "Success!",
                        text: "Data has been saved!",
                        type: "success",
                        confirmButtonColor: "#3da09a"
                    }, function() {
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

invoice.onChangeStatus = function(textFilter) {
    invoice.dataMasterInvoice([])
    var allData = invoice.dataMasterInvoiceOriginal()
    if (textFilter != "" || textFilter != undefined) {

        var Data = _.filter(allData, function(o) {
            return o.Status.indexOf(textFilter) > -1
        });
        invoice.dataMasterInvoice(Data)
    }
    invoice.renderGrid()
}

invoice.search = function(e) {
    invoice.filterAlert(true)
    if (e == "tab"){
        invoice.filterAlert(false)
    }
    invoice.textSearch(e)
    // invoice.filterindicator(true)
    if (invoice.names() == "Non Inventory") {
        invoice.getDataInvoiceNonInventory(function() {
            invoice.renderGridNonInventory()
        })
    } else {
        invoice.getDataInvoice(function() {
            invoice.renderGrid()
        })   
    }
}
invoice.printListToPdf = function(e) {
    model.Processing(true)
    var GrandTotalIDR = invoice.record.GrandTotalIDR()
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
        Id: invoice.record.ID(),
        WordGrandtotal: fixWord,
    }
    ajaxPost("/transaction/exporttopdflistinvoice", param, function(e) {
        model.Processing(false)
        var taborWindow = window.open('/res/docs/invoice/' + e, '_blank');
        taborWindow.focus();
    })
}

invoice.printListToPdfNonInventory = function(e) {
    model.Processing(true)
    var GrandTotalIDR = invoice.record.GrandTotalIDR()
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
        Id: invoice.record.ID(),
        WordGrandtotal: fixWord,
    }
    ajaxPost("/transaction/exporttopdflistinvoicenoninv", param, function(e) {
        model.Processing(false)
        var taborWindow = window.open('/res/docs/invoice/' + e, '_blank');
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
}

invoice.showEditor = function(data) {
    $('.nav-tabs a[href="#createinvoice"]').tab('show');
    if (data === undefined) {
        ko.mapping.fromJS(invoice.newRecord(), invoice.record);
        if (invoice.record.Id) invoice.record.Id("");
        if (invoice.record.ID) invoice.record.ID("");
        $('#my-checkbox').bootstrapSwitch('state', false);
        $('#checkvat-inv-cmi').bootstrapSwitch('state', false);
    } else {
        ko.mapping.fromJS(data, invoice.record);
        if (invoice.record.Id) invoice.record.Id(data.Id);
        if (invoice.record.ID) invoice.record.ID(data.Id);
        $('#my-checkbox').bootstrapSwitch('state', parseFloat(data.VAT) > 0);
        $('#checkvat-inv-cmi').bootstrapSwitch('state', data.INVCMI == true);
    }
}

invoice.editDraft = function(e) {
    // invoice.hideFilter()
    var data = []
    if (invoice.names() == "Non Inventory") {
        setTimeout(function() {
            $("th.invhideNon").css("display", "none");
            $("td.invhideNon").css("display", "none");
        }, 500);
        
        data = _.find(invoice.dataMasterInvoiceOriginalNonInventory(), function(o) {
            return o.Id == e;
        });
    } else {
        $(".invhideNon").show()
        $(".codeitem").prop('disabled', true);
        data = _.find(invoice.dataMasterInvoiceOriginal(), function(o) {
            return o.Id == e;
        });
    }
    invoice.getDataInventory(data.StoreLocationId)
    data.DateCreated = moment(data.DateCreated).format('DD-MMM-YYYY HH:mm')
    $('#dateinvoice').data('kendoDatePicker').value(new Date(data.DateCreated))
    $("#categoryDropdown").data("kendoDropDownList").value(data.Category)
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
    invoice.enableView()
    invoice.showEditor(data);

    invoice.maskingMoney()
    invoice.text("Edit Invoice")
    invoice.showCreate(false)
    invoice.showEdit(true)
    invoice.backToList(true)
    invoice.printPDFListView(false)
    invoice.checkdata()

    setTimeout(function() {
        newArr = []
        _.each(invoice.record.ListItem(), function(v, i) {
            newArr.push(invoice.listDetail(ko.mapping.toJS(v)))
        })
        invoice.record.ListItem(newArr)

        invoice.maskingMoney()
    }, 100);
}
invoice.viewDraft = function(e) {
    invoice.editDraft(e);
    invoice.printPDFListView(true)
    invoice.showEdit(false)
    invoice.disableView()
    setTimeout(function() {
        invoice.disableView()
        if (invoice.names() == "Non Inventory") {
            $(".invhideNon").hide()
            $(".codeitem").prop('disabled', false);
        } else {
            $(".invhideNon").show()
            $(".codeitem").prop('disabled', true);
        }
    }, 110);
}

invoice.viewDraft_old = function(e) {
    console.log("View Draft is render")
    invoice.hideFilter()
    var allData = invoice.dataMasterInvoiceOriginal()
    var data = _.find(allData, function(o) {
        return o.Id == e;
    });
    invoice.getDataInventory(data.StoreLocationId)
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    ko.mapping.fromJS(data, invoice.record)
    invoice.record.ID(data.Id)
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
    ko.mapping.fromJS(data, invoice.record)
    invoice.maskingMoney()
    invoice.text("View Invoice")
    invoice.showCreate(false)
    invoice.showEdit(false)
    invoice.enableView()
    if (data.VAT > 0) {
        $('#my-checkbox').bootstrapSwitch('state', true);
    } else {
        $('#my-checkbox').bootstrapSwitch('state', false);
    }
    $('#checkvat-inv-cmi').bootstrapSwitch('state', data.INVCMI == true);
    invoice.backToList(true)
    invoice.disableView()
    invoice.printPDFListView(true)
    setTimeout(function() {
        _.each(data.ListItem, function(v, i) {
            $("#itemcode_" + i).data("kendoDropDownList").value(v.CodeItem)
            drop = $("#itemcode_" + i).data("kendoDropDownList")
            drop.enable(false)
        })
    }, 100)
    $('.nav-tabs a[href="#createinvoice"]').tab('show');
}

invoice.editDraft_old = function(e) {
    invoice.hideFilter()
    var allData = invoice.dataMasterInvoiceOriginal()
    var data = _.find(allData, function(o) {
        return o.Id == e;
    });
    invoice.getDataInventory(data.StoreLocationId)
    data.DateCreated = moment(data.DateCreated).format('DD-MMM-YYYY HH:mm')
    $('#dateinvoice').data('kendoDatePicker').value(new Date(data.DateCreated))
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
    ko.mapping.fromJS(data, invoice.record)
    if (data.VAT > 0) {
        $('#my-checkbox').bootstrapSwitch('state', true);
    } else {
        $('#my-checkbox').bootstrapSwitch('state', false);
    }
    $('#checkvat-inv-cmi').bootstrapSwitch('state', data.INVCMI == true);
    invoice.record.ID(e)
    newArr = []
    _.each(invoice.record.ListItem(), function(v, i) {
        newArr.push(invoice.listDetail(ko.mapping.toJS(v)))
    })
    invoice.record.ListItem(newArr)
    invoice.maskingMoney()
    invoice.text("Edit Invoice")
    invoice.showCreate(false)
    invoice.showEdit(true)
    invoice.backToList(true)
    invoice.printPDFListView(false)
    invoice.enableView()
    invoice.checkdata()
    setTimeout(function() {
        _.each(newArr, function(v, i) {
            $("#itemcode_" + i).data("kendoDropDownList").value(v.CodeItem())
        })
    }, 100)
    $('.nav-tabs a[href="#createinvoice"]').tab('show');
}

invoice.backList = function() {
    $('.nav-tabs a[href="#listinvoice"]').tab('show')
    $("#listinvoice").addClass("active");
    $("#createinvoice").removeClass("active");
}

invoice.disableView = function() {
    $(".formInput").attr("disabled", "disabled")
    $(".invhide").hide()
    // $(".btnDeleteSummary").attr("disabled", "disabled")
    // $("#addnewitem").attr("disabled", "disabled")
    $("#my-checkbox").bootstrapSwitch('disabled', true)
    $("#checkvat-inv-cmi").bootstrapSwitch('disabled', true)
    $('#dateinvoice').data('kendoDatePicker').enable(false);
    var dropDown2 = $("#customercode").data("kendoDropDownList");
    dropDown2.enable(false);
    var dropDown3 = $("#storelocation").data("kendoDropDownList");
    dropDown3.enable(false);
    var dropDown4 = $("#categoryDropdown").data("kendoDropDownList");
    dropDown4.enable(false);
    $("#salescode").data("kendoDropDownList").enable(false)
    $("select.codeitem").each(function(idx, ele) {
        $(this).data("kendoDropDownList").enable(false);
    });
    // var dropDown3 = $("#accountnumber").data("kendoDropDownList");
    // dropDown3.enable(false);
}

invoice.enableView = function() {
    $(".formInput").removeAttr("disabled")
    $(".invhide").show()
    // $(".btnDeleteSummary").removeAttr("disabled")
    // $("#addnewitem").removeAttr("disabled")
    $("#my-checkbox").bootstrapSwitch('disabled', false)
    $("#checkvat-inv-cmi").bootstrapSwitch('disabled', false)
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#dateinvoice').data('kendoDatePicker').enable(condition);
    var dropDown2 = $("#customercode").data("kendoDropDownList");
    dropDown2.enable(true);
    var dropDown3 = $("#storelocation").data("kendoDropDownList");
    dropDown3.enable(true);
    var dropDown4 = $("#categoryDropdown").data("kendoDropDownList");
    dropDown4.enable(true);
    $("#salescode").data("kendoDropDownList").enable(true)
    // var dropDown3 = $("#accountnumber").data("kendoDropDownList");
    // dropDown3.enable(true);
    $("#my-checkbox").bootstrapSwitch({
        disabled: false
    });
}

invoice.resetView = function() {
    ko.mapping.fromJS(invoice.newRecord(), invoice.record)
    invoice.record.Rate(1)
    invoice.record.DateCreated(moment(new Date()).format('DD-MMM-YYYY'))
    $(".formInput").val("")
    $("#customername").val("")
    $(".Amount").val("")
    $("#checkvat-inv-cmi").bootstrapSwitch('state', false)
    invoice.text("Create Invoice")
    $('#customercode').data('kendoDropDownList').value(-1);
    // $('#storelocation').data('kendoDropDownList').value(-1);
    // $('#accountnumber').data('kendoDropDownList').value(-1);
    invoice.enableView()
    $("#my-checkbox").bootstrapSwitch('state', false);
    invoice.dataDropDownInventory([])
}

invoice.hideFilter = function() {
    var panelFilter = $('.panel-filter');
    var panelContent = $('.panel-content');
    panelFilter.hide();
    panelContent.attr('class', 'col-md-12 col-sm-12 ez panel-content');
    $('.breakdown-filter').removeAttr('style');
    invoice.TitelFilter(" Show Filter");
}

invoice.setDate = function() {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
invoice.fromPOINVSummary = function() {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    if (num != null) {
        var allData = invoice.dataMasterInvoiceOriginal()
        var data = _.find(allData, function(o) {
            return o.DocumentNo == num;
        });
        // console.log(num, data)
        if (data != undefined) {
            invoice.viewDraft(data.Id)
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning"
            }, function() {
                window.location.assign("/transaction/invoice")
            });
        }
    }
}
invoice.processsPlusPlus = function() {
    processs += 1
    // console.log(processs)
    if (processs >= 4) {
        setTimeout(function() {
            invoice.fromPOINVSummary()
        }, 300)
        // setTimeut(purchasepayment.createPPFromPOInvSummary(),10000) 
    }
}

invoice.getDataInventory = function(val) {
    var param = {
        LocationId: val
    }
    model.Processing(true)

    ajaxPost('/master/getdatainventory', param, function(res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        invoice.dataMasterInventory(res.Data)
        var DataInventory = res.Data
        for (i in DataInventory) {
            DataInventory[i].Kode = DataInventory[i].INVID + ""
            DataInventory[i].Name = DataInventory[i].INVID + " - " + DataInventory[i].INVDesc
        }
        invoice.dataDropDownInventory(DataInventory)

        model.Processing(false)
    })
}

invoice.onChangeCodeItem = function(value, index) {
    findaccount = _.find(invoice.dataMasterInventory(), {
        INVID: value
    })
    // console.log(findaccount.INVDesc)
    invoice.record.ListItem()[index].Item(findaccount.INVDesc);
}

invoice.onSelect = function(e) {
    var dataItem = this.dataItem(e.item);
    //    console.log(dataItem.SalesName)
    invoice.record.SalesName(dataItem.SalesName)
};

invoice.addCategory = function(id, val) {
    var widget = $("#" + id).getKendoDropDownList();
     var dataSource = widget.dataSource;
    swal({
        title: "Are you sure?",
        text: "You want add new Department",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: false,
        closeOnCancel: false
    }, function(isConfirm) {
        if(isConfirm) {
            dataSource.add({
                _id: "",
                code: "",
                name: val
            });
            dataSource.one("sync", function () {
                widget.select(dataSource.view().length - 1);
            });
            dataSource.sync();
            var dropDownList = $("#categoryDropdown").data("kendoDropDownList")
            dropDownList.value(val)
            dropDownList.trigger("change")
            model.Processing(true)
            var url = "/master/insertnewinvoicecategory"
            var param = {
                Name: val
            }
            ajaxPost(url, param, function(res) {
                model.Processing(false)
                swal({
                    title: "Success!",
                    text: "Success Add New Data",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                })
            })
        }else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}

invoice.dropdownCategory = function() {
    $("#categoryDropdown").html("")
    ajaxPost("/master/getdatainvoicecategory", {}, function(res) {
        var result = _.sortBy(res.Data, [function(o) { return o.categoryname; }]);
        $("#categoryDropdown").kendoDropDownList({
            filter: "contains",
            dataTextField: "categoryname",
            dataValueField: "categoryname",
            dataSource: result,
            optionLabel: "Select ..",
            noDataTemplate: $("#noDataTemplate").html(),
            change: function(e) {
                var dataItem = this.dataItem();
                console.log(dataItem, this)
                invoice.dataMasterCategory(dataItem.name)
            }
        })
        var dropDownList = $("#categoryDropdown").data("kendoDropDownList");
        if (dropDownList.value() != "") {
            dropDownList.value("");
            dropDownList.trigger("change");
        }
    })
}
invoice.getDataStoreLocation = function() {
    model.Processing(true)
    ajaxPost("/master/getdatalocationbyuser", {}, function(res) {
        model.Processing(false)
        if (res.IsError) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        _.each(res.Data, function(e) {
            e.Text = e.LocationID + "-" + e.LocationName
            e.Value = e.LocationID.toString()
        })
        invoice.dataMasterLocation(res.Data)
    })
}
invoice.onChangeLocation = function(value) {
    var code = parseInt(value)
    var data = _.find(invoice.dataMasterLocation(), function(r) {
        return r.LocationID == code
    })
    invoice.record.StoreLocationName(data.LocationName)
    invoice.getDataInventory(code)
}
invoice.choose = function(choosenote) {
    $('.nav-tabs a[href="#List"]').tab('show')
    invoice.textSearch("")
    $("#textSearch").val("")
    invoice.filterindicator(false)
    if (choosenote == "ACTIVA") {
        invoice.names("Non Inventory")
        setTimeout(function() {
            $(".invhideNon").hide()
            $(".typehide").show()
        }, 50);
        $("#noninv").css({
            "background-color": "#ffffff",
            "color": "#3ea49d",
            "border-bottom": "4px solid #f4222d"
        });
        $("#inv").css({
            "background-color": "#45b6af",
            "color": "white",
            "border-color": "#3ea49d",
            "border-bottom": "1px solid #3ea49d"
        });

        $(".item").prop('disabled', false);

        invoice.getDataInvoiceNonInventory(function() {
            invoice.renderGridNonInventory()
        })

    } else {
        invoice.getDataInvoice(function() {
            invoice.renderGrid()
        })
        invoice.names("Inventory")
        setTimeout(function() {
            $(".invhideNon").show()
            $(".typehide").hide()
        }, 50);

        $(".item").prop('disabled', true);
        $("#inv").css({
            "background-color": "#ffffff",
            "color": "#3ea49d",
            "border-bottom": "4px solid #f4222d"
        });
        $("#noninv").css({
            "background-color": "#45b6af",
            "color": "white",
            "border-color": "#3ea49d",
            "border-bottom": "1px solid #3ea49d"
        });

    }
    // invoice.resetView()
}
invoice.init = function() {
    var now = new Date()

    invoice.date(moment(now).format("DD MMMM YYYY"))
    // invoice.initDate()
    // invoice.setDate()
    invoice.getMasterLocation()
    invoice.dropdownCategory()
    invoice.getDataCustomer()
    invoice.getDataPurchaseOrder()
    invoice.getDataAccount()
    invoice.getDataInvoiceNonInventory(function() {
        invoice.renderGridNonInventory()
    })
    invoice.text("Create Invoice")
    invoice.maskingMoney()
    invoice.switchButton()
    invoice.getDateNow()
    invoice.getDataInventory(0)
    invoice.getDataSales()
    invoice.getDataStoreLocation()
    invoice.names("Non Inventory")
    $("#noninv").css({
        "background-color": "#ffffff",
        "color": "#3ea49d",
        "border-bottom": "4px solid #f4222d"
    });

}

invoice.detailReportPdf = function() {
    var valuefilter = invoice.textCustomerSearch()
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: moment(startdate).format("YYYY-MM-DD"),
        DateEnd: moment(enddate).format("YYYY-MM-DD"),
        ReportType: "summary",
        ReportBy: "Customer",
        ValueFilter: valuefilter
    }
    model.Processing(true)
    ajaxPost("/report/detailreportpdf", param, function(res) {
        if (res.IsError) {
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}

invoice.detailReportPdfNonInv = function() {
    var valuefilter = invoice.textCustomerSearch()
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: moment(startdate).format("YYYY-MM-DD"),
        DateEnd: moment(enddate).format("YYYY-MM-DD"),
        ReportType: "summary",
        ReportBy: "Customer",
        ValueFilter: valuefilter
    }
    model.Processing(true)
    ajaxPost("/transaction/detailreportpdfnoninv", param, function(res) {
        if (res.IsError) {
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}

invoice.saveSwitch = function() {
    var validator = $("#createinvoice").kendoValidator().data("kendoValidator")
    // console.log(validator.validate())
    if (validator.validate()) {
        if (invoice.names() == "Non Inventory") {
            invoice.saveDataNonInventory()
        } else {
            invoice.saveData()
        }
    }
}

invoice.postingSwitch = function() {
    if (invoice.names() == "Non Inventory") {
        invoice.PostingDataNonInventory()
    } else {
        invoice.PostingData()
    }
}

invoice.detailReportPdfSwitch = function() {
    if (invoice.names() == "Non Inventory") {
        invoice.detailReportPdfNonInv()
    } else {
        invoice.detailReportPdf()
    }
}

invoice.printSwitch = function() {
    if (invoice.names() == "Non Inventory") {
        invoice.printListToPdfNonInventory()
    } else {
        invoice.printListToPdf()
    }
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
        digitsGroup[l].filter(function(e) {
            return e !== null
        });
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

$(function() {
    invoice.init();
    
    $("#textSearch").on("keyup blur change", function () {
        invoice.filterText();
    });
})

invoice.onChangeDateStart = function(val){
    if (val.getTime()>invoice.DateEnd().getTime()){
        invoice.DateEnd(val)
    }
}

// invoice.initDate = function () {
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

invoice.getMasterLocation= function() {
    $.ajax({
        url: "/transferorder/getuserlocations",
        success: function(data) {
            invoice.warehouse.removeAll();
            $(data).each(function(ix, ele) {
                invoice.dataLocation.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
                invoice.warehouse.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
            });
            invoice.valueStorehouse(data[0]["LocationID"])
        }
    });
}

invoice.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["AccountName", "DocumentNo", "CustomerName", "CustomerCode", "PoNumber", "SalesCode", "SalesName", "StoreLocationName", "Status"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridListInvoice").data("kendoGrid").dataSource.query({ filter: filter });
}