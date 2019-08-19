model.Processing(false)

var purchasecreditmemo = {}
purchasecreditmemo.dataMasterSupplier = ko.observableArray([])
purchasecreditmemo.dataDropDownSupplier = ko.observableArray([])
purchasecreditmemo.dataDropDownSupplierFilter = ko.observableArray([])
purchasecreditmemo.dataMasterAccount = ko.observableArray([])
purchasecreditmemo.dataDropDownAccount = ko.observableArray([])
purchasecreditmemo.dataMasterpurchasecreditmemo = ko.observableArray([])
purchasecreditmemo.dataListpurchasecreditmemo = ko.observableArray([])
purchasecreditmemo.dataListDraftpurchasecreditmemo = ko.observableArray([])
purchasecreditmemo.dataMasterpurchasecreditmemoOriginal = ko.observableArray([])
purchasecreditmemo.dataMasterPurchaseInvoice = ko.observableArray([])
purchasecreditmemo.type = ko.observableArray([])
purchasecreditmemo.checkedvat = ko.observable(false)
purchasecreditmemo.text = ko.observable()
purchasecreditmemo.textSearch = ko.observable()
purchasecreditmemo.statustext = ko.observable()
purchasecreditmemo.sequenceNumber = ko.observable()
purchasecreditmemo.showCreate = ko.observable(true)
purchasecreditmemo.backToDraft = ko.observable(false)
purchasecreditmemo.postDraft = ko.observable(false)
purchasecreditmemo.deleteDraft = ko.observable(false)
purchasecreditmemo.printPDF = ko.observable(false)
purchasecreditmemo.reset = ko.observable(false)
purchasecreditmemo.setDocumentNumber = ko.observable('')
purchasecreditmemo.rowSpan = ko.observable(5)
purchasecreditmemo.DatePageBar = ko.observable()
purchasecreditmemo.BoolVat = ko.observable(false)
purchasecreditmemo.textSupplierSearch = ko.observable()
purchasecreditmemo.filterindicator = ko.observable(false)
purchasecreditmemo.TotalAll = ko.observable(0)
purchasecreditmemo.GrandTotalAll = ko.observable(0)
purchasecreditmemo.dataDropDownSales = ko.observableArray([])
purchasecreditmemo.dataMasterSales = ko.observableArray([])

purchasecreditmemo.payment = [{
    value: "CASH",
    text: "Cash"
}, {
    value: "INSTALLMENT",
    text: "Installment"
}]

purchasecreditmemo.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    purchasecreditmemo.DatePageBar(page)
}
purchasecreditmemo.newRecord = function() {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumber: "",
        SupplierCode: "",
        SupplierName: "",
        AccountCode: 0,
        Payment: "",
        // Type: "",
        ListDetail: [],
        Status: "",
        TotalIDR: 0,
        TotalUSD: 0,
        Discount: 0,
        VAT: 0,
        GrandTotalIDR: 0,
        GrandTotalUSD: 0,
        Remark: "",
        Rate: 0,
        Currency: "",
        DownPayment: 0,
        Department: "",
        SalesCode: "",
        SalesName: "",  
    }
    page.ListDetail.push(purchasecreditmemo.listDetail({}))
    return page
}

purchasecreditmemo.listDetail = function(data) {
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
        //Menghitung Total Amount IDR
        total = FormatCurrency(e) * FormatCurrency(x.PriceIDR())
        x.AmountIDR(total)

        //Menghitung Total Amount USD
        totalUSD = FormatCurrency(e) * FormatCurrency(x.PriceUSD())
        x.AmountUSD(totalUSD)


        //Menghitung Total
        var data = ko.mapping.toJS(purchasecreditmemo.record)
        if (purchasecreditmemo.record.ListDetail()[0].PriceIDR > 0 || total > 0) {
            var alldataidr = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
            var totalIDR = _.sumBy(alldataidr, function(v) {
                return v.AmountIDR
            })
            purchasecreditmemo.record.TotalIDR(totalIDR)
            purchasecreditmemo.TotalAll(ChangeToRupiah(purchasecreditmemo.record.TotalIDR()))

        } else {
            var alldatausd = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
            var totalUSD = _.sumBy(alldatausd, function(v) {
                return v.AmountUSD
            })
            purchasecreditmemo.record.TotalUSD(totalUSD)
            var totalIDRRATE = totalUSD * FormatCurrency(purchasecreditmemo.record.Rate())
            purchasecreditmemo.record.TotalIDR(totalIDRRATE)
            purchasecreditmemo.TotalAll(ChangeToRupiah(purchasecreditmemo.record.TotalUSD()))

        }

    })

    x.PriceIDR.subscribe(function(e) {
        x.AmountIDR(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var alldataidr = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
        var totalIDR = _.sumBy(alldataidr, function(v) {
            return v.AmountIDR
        })
        // purchasecreditmemo.record.Total(ChangeToRupiah(totalIDR))
        purchasecreditmemo.record.TotalIDR(totalIDR)
        purchasecreditmemo.record.Currency("IDR")
        purchasecreditmemo.TotalAll(ChangeToRupiah(purchasecreditmemo.record.TotalIDR()))
        purchasecreditmemo.checkdata()

    })

    x.PriceUSD.subscribe(function(e) {
        x.AmountUSD(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var data = ko.mapping.toJS(purchasecreditmemo.record)
        // x.PriceIDR(FormatCurrency(e)*FormatCurrency(data.Rate))
        var alldatausd = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
        var totalUSD = _.sumBy(alldatausd, function(v) {
            return v.AmountUSD
        })
        // purchasecreditmemo.record.Total(ChangeToRupiah(totalIDR))
        purchasecreditmemo.record.TotalUSD(totalUSD)
        var totalIDRRATE = totalUSD * FormatCurrency(purchasecreditmemo.record.Rate())
        purchasecreditmemo.record.TotalIDR(totalIDRRATE)
        purchasecreditmemo.record.Currency("USD")
        purchasecreditmemo.TotalAll(ChangeToRupiah(purchasecreditmemo.record.TotalUSD()))
        purchasecreditmemo.checkdata()
    })

    return x
}

purchasecreditmemo.record = ko.mapping.fromJS(purchasecreditmemo.newRecord())

purchasecreditmemo.switchButton = function() {
    $('#checkvat').on('switchChange.bootstrapSwitch', function(event, state) {
        var data = ko.mapping.toJS(purchasecreditmemo.record)
        purchasecreditmemo.BoolVat(state)
        if (state) {
            if (purchasecreditmemo.record.Currency() == "IDR") {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalIDR)
                purchasecreditmemo.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
                var Totalall = GT + totVat
                purchasecreditmemo.record.GrandTotalIDR(Totalall)
                purchasecreditmemo.GrandTotalAll(ChangeToRupiah(Totalall))

            } else {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalUSD)
                purchasecreditmemo.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
                var Totalall = GT + totVat
                purchasecreditmemo.record.GrandTotalUSD(Totalall)
                purchasecreditmemo.record.GrandTotalIDR(Totalall * FormatCurrency(data.Rate))
                purchasecreditmemo.GrandTotalAll(ChangeToRupiah(Totalall))

            }
        } else {
            purchasecreditmemo.record.VAT(ChangeToRupiah(0))
            if (purchasecreditmemo.record.Currency() == "IDR") {
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
                purchasecreditmemo.record.GrandTotalIDR(GT)
                purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT))

            } else {
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
                purchasecreditmemo.record.GrandTotalUSD(GT)
                purchasecreditmemo.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
                purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT))

            }
        }
    });
}

purchasecreditmemo.TotalAll.subscribe(function(e) {
    var data = ko.mapping.toJS(purchasecreditmemo.record)
    if (purchasecreditmemo.BoolVat()) {
        var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(e)
        var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(e)
        //console.log(totVat)
        purchasecreditmemo.record.VAT(totVat)
        if (purchasecreditmemo.record.Currency() == "IDR") {
            purchasecreditmemo.record.GrandTotalIDR(GT + totVat)
            purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT + totVat))
        } else {
            purchasecreditmemo.record.GrandTotalUSD(GT + totVat)
            purchasecreditmemo.record.GrandTotalIDR((GT + totVat) * FormatCurrency(data.Rate))
            purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT + totVat))
        }
    } else {
        purchasecreditmemo.record.VAT(0)
        if (data.ListDetail.length != 0) {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(e)
            if (purchasecreditmemo.record.Currency() == "IDR") {
                purchasecreditmemo.record.GrandTotalIDR(GT)
                purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT))
            } else {
                purchasecreditmemo.record.GrandTotalUSD(GT)
                purchasecreditmemo.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
                purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT))
            }

        } else {
            purchasecreditmemo.GrandTotalAll(0)
        }
    }
})

purchasecreditmemo.record.Discount.subscribe(function(e) {
    var data = ko.mapping.toJS(purchasecreditmemo.record)
    if (FormatCurrency(e) > 100 || FormatCurrency(e) < 0) {
        return swal({
            title: "Warning!",
            text: "your discount are irational",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function() {
            purchasecreditmemo.record.Discount(0)
        });
    }


    if (purchasecreditmemo.record.Currency() == "IDR") {
        var GT = ((100 - FormatCurrency(e)) / 100) * FormatCurrency(data.TotalIDR)
        var totVat = 0
        if (purchasecreditmemo.BoolVat()) {
            totVat = ((100 - FormatCurrency(e)) / 1000) * FormatCurrency(data.TotalIDR)
            purchasecreditmemo.record.VAT(ChangeToRupiah(totVat))
        } else {
            purchasecreditmemo.record.VAT(ChangeToRupiah(0))
            totVat = 0
        }
        purchasecreditmemo.record.GrandTotalIDR(GT + totVat)
        purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT + totVat))
    } else {
        var GT = ((100 - FormatCurrency(e)) / 100) * FormatCurrency(data.TotalUSD)
        var totVat = 0
        if (purchasecreditmemo.BoolVat()) {
            var totVat = ((100 - FormatCurrency(e)) / 1000) * FormatCurrency(data.TotalUSD)
            purchasecreditmemo.record.VAT(ChangeToRupiah(totVat))
        } else {
            purchasecreditmemo.record.VAT(ChangeToRupiah(0))
            totVat = 0
        }

        purchasecreditmemo.record.GrandTotalUSD(GT + totVat)
        purchasecreditmemo.record.GrandTotalIDR((GT + totVat) * FormatCurrency(data.Rate))
        purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT + totVat))

    }

})

purchasecreditmemo.record.VAT.subscribe(function(e) {
    var data = ko.mapping.toJS(purchasecreditmemo.record)
    // purchasecreditmemo.BoolVat(state)
    if (purchasecreditmemo.BoolVat()) {
        if (purchasecreditmemo.record.Currency() == "IDR") {
            var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalIDR)
            purchasecreditmemo.record.VAT(ChangeToRupiah(totVat))
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
            var Totalall = GT + totVat
            purchasecreditmemo.record.GrandTotalIDR(Totalall)
            purchasecreditmemo.GrandTotalAll(ChangeToRupiah(Totalall))

        } else {
            var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalUSD)
            purchasecreditmemo.record.VAT(ChangeToRupiah(totVat))
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
            var Totalall = GT + totVat
            purchasecreditmemo.record.GrandTotalUSD(Totalall)
            purchasecreditmemo.record.GrandTotalIDR(Totalall * FormatCurrency(data.Rate))
            purchasecreditmemo.GrandTotalAll(ChangeToRupiah(Totalall))

        }
    } else {
        purchasecreditmemo.record.VAT(ChangeToRupiah(0))
        if (purchasecreditmemo.record.Currency() == "IDR") {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
            purchasecreditmemo.record.GrandTotalUSD(GT)
            purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT))

        } else {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
            purchasecreditmemo.record.GrandTotalUSD(GT)
            purchasecreditmemo.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
            purchasecreditmemo.GrandTotalAll(ChangeToRupiah(GT))

        }
    }
})
purchasecreditmemo.record.DownPayment.subscribe(function(e) {
    //console.log(e)
    if (parseFloat(e) > 100 || parseFloat(e) < 0) {
        return swal({
            title: "Warning!",
            text: "your Down Payment are irational",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function() {
            purchasecreditmemo.record.DownPayment(0)
        });
    }
})

// purchasecreditmemo.record.Rate.subscribe(function (e) {
//     var data = ko.mapping.toJS(purchasecreditmemo.record)
//     if (data.ListDetail[0].PriceUSD > 0) {
//         for (var i = 0; i < data.ListDetail.length; i++) {
//             var alldatausd = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
//             var totalIDR = _.sumBy(alldatausd, function (v) {
//                 return v.AmountUSD * FormatCurrency(e)
//             })
//             purchasecreditmemo.record.Total(ChangeToRupiah(totalIDR))
//         }
//     }
// })

purchasecreditmemo.getDataSales = function() {
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
        purchasecreditmemo.dataMasterSales(res.Data)
        var DataSales = res.Data
        for (i in DataSales) {
            DataSales[i].Kode = DataSales[i].SalesID + ""
            DataSales[i].Name = DataSales[i].SalesID + " - " + DataSales[i].SalesName
        }
        purchasecreditmemo.dataDropDownSales(DataSales)

        model.Processing(false)
    })
}

purchasecreditmemo.getDataSupplier = function() {
    model.Processing(true)

    ajaxPost('/transaction/getsupplier', {}, function(res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchasecreditmemo.dataMasterSupplier(res.Data)
        purchasecreditmemo.dataDropDownSupplier(res.Data)
        var DataSupplier = res.Data
        for (i in DataSupplier) {
            DataSupplier[i].Kode = DataSupplier[i].Kode + ""
            DataSupplier[i].Name = DataSupplier[i].Kode + " - " + DataSupplier[i].Name
        }
        purchasecreditmemo.dataDropDownSupplierFilter(DataSupplier)

        model.Processing(false)
    })
}

purchasecreditmemo.getDataTypePurchase = function() {
    model.Processing(true)

    ajaxPost('/transaction/gettypepurchase', {}, function(res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchasecreditmemo.type(res.Data)
        model.Processing(false)
    })
}

purchasecreditmemo.getDataAccount = function() {
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
        purchasecreditmemo.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        purchasecreditmemo.dataDropDownAccount(DataAccount)

        model.Processing(false)
    })
}

purchasecreditmemo.onChangeAccountNumber = function(value) {
    var result = _.filter(purchasecreditmemo.dataMasterSupplier(), {
        'Kode': value
    })[0].Name
    purchasecreditmemo.record.SupplierName(result);
}

purchasecreditmemo.maskingMoney = function() {
    $('.currency').inputmask("numeric", {
        radixPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        // prefix: '$', //No Space, this will truncate the first character
        rightAlign: false,
        // oncleared: function () { self.Value(''); }
    });
}

purchasecreditmemo.getDatapurchasecreditmemo = function(callback) {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {}
    if (purchasecreditmemo.filterindicator() == true) {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true,
            TextSearch: purchasecreditmemo.textSearch(),
            SupplierCode: purchasecreditmemo.textSupplierSearch(),
        }
    } else {
        startdate = moment().startOf('month').format('YYYY-MM-DD hh:mm')
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: false,
        }
    }
    model.Processing(true)
    ajaxPost('/transaction/getdatapurchasecreditmemo', param, function(res) {
        if (res.IsError) {
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
        purchasecreditmemo.dataMasterpurchasecreditmemo(res.Data)
        purchasecreditmemo.dataMasterpurchasecreditmemoOriginal(res.Data)
        model.Processing(false)
        callback()
        purchasecreditmemo.fromPOINVSummary()
    }, function() {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

purchasecreditmemo.GetDataPurchaseInvoice = function() {
    model.Processing(true)
    ajaxPost('/transaction/getdatapurchaseinvoiceinventoryforcreditmemo', {}, function(res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchasecreditmemo.dataMasterPurchaseInvoice(res.Data)
        model.Processing(false)
    })
}
purchasecreditmemo.renderGrid = function() {
    var data = purchasecreditmemo.dataMasterpurchasecreditmemo();
    // for (var i = 0; i < data.le.Kondisingth; i++) {
    //     data[i].Kondisi = (data[i].User ==  )
    // }
    // if (typeof $('#gridListpurchasecreditmemo').data('kendoGrid') !== 'undefined') {
    //     $('#gridListpurchasecreditmemo').data('kendoGrid').setDataSource(new kendo.data.DataSource({
    //         data: data,
    //     }))
    //     return
    // }

    var columns = [

        {
            title: 'Action',
            width: 50,
            template: "# if (User == userinfo.usernameh() || userinfo.usernameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' ) {#<button onclick='purchasecreditmemo.viewDraft(\"#: _id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}else{#<button onclick='purchasecreditmemo.viewDraft(\"#: _id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}#",
        }, {
            field: 'DateStr',
            title: 'Date',
            width: 160,
        }, {
            field: 'DocumentNumber',
            title: 'Order #',
            width: 160,
        }, {
            field: 'SupplierName',
            title: 'Supplier Name',
            width: 200,
        }, {
            field: 'GrandTotalIDR',
            title: 'Order Total ',
            width: 200,
            // template: "#=ChangeToRupiah(GrandTotalUSD)#"
            attributes: {
                "class": "rightAction",
            },
            template: "#if ( Currency == 'USD') {# $ #: ChangeToRupiah(GrandTotalUSD) # #} else {# Rp. #: ChangeToRupiah(GrandTotalIDR) # #}#"

        }, {
            field: 'Remark',
            title: 'Remark',
            width: 200,
        }
    ]

    $('#gridListPurchasecreditmemo').kendoGrid({
        dataSource: {
            data: data,
            sort: {
                field: 'DatePosting',
                dir: 'desc',
            }
        },
        height: 500,
        width: 140,
        sortable: true,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "PurchaseCreditMemo", function(row, sheet){
                console.log(row.cells.length)
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
                        if (ci == 0) {
                            cell.format = "dd-MM-yyyy"; 
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
purchasecreditmemo.exportExcel = function () {
    $("#gridListPurchasecreditmemo").getKendoGrid().saveAsExcel();
}

purchasecreditmemo.checkdata = function() {
    if (FormatCurrency(purchasecreditmemo.record.ListDetail()[0].PriceUSD()) == 0 && FormatCurrency(purchasecreditmemo.record.ListDetail()[0].PriceIDR()) == 0) {
        $(".priceidr").removeAttr("disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(purchasecreditmemo.record.ListDetail()[0].PriceUSD()) > 0) {
        $(".priceidr").attr("disabled", "disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(purchasecreditmemo.record.ListDetail()[0].PriceUSD()) == 0) {
        $(".priceusd").attr("disabled", "disabled");
        $(".priceidr").removeAttr("disabled");
    }
}

purchasecreditmemo.addNewItem = function() {
    purchasecreditmemo.record.ListDetail.push(purchasecreditmemo.listDetail({}))
    purchasecreditmemo.checkdata()
    purchasecreditmemo.maskingMoney()
}

purchasecreditmemo.setDate = function() {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}


purchasecreditmemo.removeRow = function() {
    purchasecreditmemo.checkdata()
    purchasecreditmemo.record.ListDetail.remove(this)
    if (purchasecreditmemo.record.ListDetail().length == 0) {
        purchasecreditmemo.record.ListDetail.push(purchasecreditmemo.listDetail({}))
    }
    //Menghitung Total 
    var alldataidr = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
    var totalIDR = _.sumBy(alldataidr, function(v) {
        return v.AmountIDR
    })

    var alldata = ko.mapping.toJS(purchasecreditmemo.record.ListDetail())
    var totalUSD = _.sumBy(alldata, function(v) {
        return v.AmountUSD
    })

    if (FormatCurrency(purchasecreditmemo.record.ListDetail()[0].PriceIDR()) > 0) {
        purchasecreditmemo.record.Total(ChangeToRupiah(totalIDR))
    }

    if (FormatCurrency(purchasecreditmemo.record.ListDetail()[0].PriceUSD()) > 0) {
        purchasecreditmemo.record.Total(ChangeToRupiah(totalUSD))
    }

    var data = ko.mapping.toJS(purchasecreditmemo.record)
    //Menghitung VAT           
    var nomvat = ((100 - data.Discount) / 1000) * FormatCurrency(data.Total)
    purchasecreditmemo.record.VAT(ChangeToRupiah(nomvat))

    // Menghitung GrandTotal
    var totalnumber = FormatCurrency(data.Total)
    var nomgrandtotal = (((100 - data.Discount) / 100) * totalnumber) + nomvat
    purchasecreditmemo.record.GrandTotal(ChangeToRupiah(nomgrandtotal))
}

purchasecreditmemo.getDocNumberPO = function() {
    model.Processing(true)
    ajaxPost('/transaction/getlastnumberpo', {}, function(res) {
        var tgl = moment(Date()).format('/DDMMYY/')
        var zr
        if (res.Number < 10) {
            zr = '00'
        } else if (res.Number >= 10 && res.Number < 100) {
            zr = '0'
        } else {
            zr = ''
        }

        purchasecreditmemo.record.DocumentNumber('PO' + tgl + zr + res.Number);
        purchasecreditmemo.sequenceNumber(res.Number)
        purchasecreditmemo.setDocumentNumber('PO' + tgl + zr + res.Number);
    })
}

purchasecreditmemo.saveData = function() {

    var change = ko.mapping.toJS(purchasecreditmemo.record)
    if (change.SupplierName == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't choose the supplier code",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (change.Payment == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't choose the payment",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    // if (change.Type == "") {
    //     return swal({
    //         title: 'Warning!',
    //         text: "You haven't choose the type",
    //         type: "info",
    //         confirmButtonColor: "#3da09a"
    //     })
    // }
    //  if (FormatCurrency(change.Rate) == 0) {
    //      return swal({
    //          title: 'Warning!',
    //          text: "You haven't input the rate",
    //          type: "info",
    //          confirmButtonColor: "#3da09a"
    //      })
    //  }
    // if (change.DownPayment == "") {
    //     return swal({
    //         title: 'Warning!',
    //         text: "You haven't input the Down Payment",
    //         type: "info",
    //         confirmButtonColor: "#3da09a"
    //     })
    // }
    if (change.Remark == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't input the remark",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    var itmCount = 0;
    for (var i = 0; i < change.ListDetail.length; i++) {
        change.ListDetail[i].Qty = parseInt(change.ListDetail[i].Qty)
        change.ListDetail[i].PriceIDR = FormatCurrency(change.ListDetail[i].PriceIDR)
        change.ListDetail[i].AmountIDR = FormatCurrency(change.ListDetail[i].PriceIDR * change.ListDetail[i].Qty)

        //in case user use IDR only
        if (FormatCurrency(change.ListDetail[i].PriceUSD) == 0) {
            var USD = FormatCurrency(change.ListDetail[i].PriceIDR) / FormatCurrency(change.Rate)
            change.ListDetail[i].PriceUSD = USD
            change.ListDetail[i].AmountUSD = USD * FormatCurrency(change.ListDetail[i].Qty)
            change.Currency = "IDR"
        } else {
            change.Currency = "USD"
            change.ListDetail[i].PriceUSD = FormatCurrency(change.ListDetail[i].PriceUSD)
            change.ListDetail[i].AmountUSD = FormatCurrency(change.ListDetail[i].PriceUSD * change.ListDetail[i].Qty)
        }
        if (change.ListDetail[i].Qty != 0)
            itmCount++;
    }
    if (itmCount == 0) {
        return swal({
            title: 'Warning!',
            text: "You have not inputted any item",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    change.DatePosting = $('#datepurchase').data('kendoDatePicker').value()
    change.DateStr = moment($('#datepurchase').data('kendoDatePicker').value()).format("DD-MMM-YYYY");
    change.Status = "PO"
    // change.StatusPayment = "open"
    change.Rate = FormatCurrency(change.Rate)
    change.Discount = FormatCurrency(change.Discount)
    change.DownPayment = FormatCurrency(change.DownPayment)
    change.VAT = FormatCurrency(change.VAT)
    change.Remark = change.Remark
    if (change.Currency == "IDR" && change.Rate != 0 && change.Discount == 0) {
        change.GrandTotalIDR = change.TotalIDR + change.VAT
        change.GrandTotalUSD = change.GrandTotalIDR / change.Rate
    }
    var param = {
        Data: change,
        // LastNumber: purchasecreditmemo.sequenceNumber(),
    }
    var url = "/transaction/insertdraftpurchasecreditmemo"
    swal({
        title: "Are you sure?",
        text: "You will submit this Purchase Credit Memo",
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
                        window.location.assign("/transaction/purchasecreditmemo")
                        $("#btnSave").attr("disabled", "disabled")
                    });
                    purchasecreditmemo.getDatapurchasecreditmemo(function() {
                        purchasecreditmemo.renderGrid()
                    })
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

purchasecreditmemo.printToPdf = function() {
    model.Processing(true)
    var param = {
        Id: purchasecreditmemo.record.ID(),
    }
    ajaxPost("/transaction/exporttopdfpurchasecreditmemo", param, function(e) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/purchasecreditmemo/' + e, '_blank');
        tabOrWindow.focus();
        // console.log(e);
        // window.open(e);
        // var pom = document.createElement('a');
        // pom.setAttribute('href', "/res/docs/purchasecreditmemo/" + e);
        // pom.setAttribute('download', e);
        // pom.click();

        // console.log(pom)
        /*var url = ("{{BaseUrl}}/res/docs/purchasecreditmemo/" + e)
        var link = document.createElement('a');
        link.href = url;
        link.dispatchEvent(new MouseEvent('click'));*/
        /*var fileURL = '{{BaseUrl}}res/docs/purchasecreditmemo/' + e
        var fileName = e
        if (!window.ActiveXObject) {
            var save = document.createElement('a');
            save.href = fileURL;
            save.target = '_blank';
            save.download = fileName || 'unknown';

            var evt = new MouseEvent('click', {
                'view': window,
                'bubbles': true,
                'cancelable': false
            });
            save.dispatchEvent(evt);

            //console.log(save)

            (window.URL || window.webkitURL).revokeObjectURL(save.href);
        }

        // for IE < 11
        else if ( !! window.ActiveXObject && document.execCommand)     {
            var _window = window.open(fileURL, '_blank');
            _window.document.close();
            _window.document.execCommand('SaveAs', true, fileName || fileURL)
            _window.close();
        }*/
    })
}

purchasecreditmemo.viewDraft = function(e) {
    purchasecreditmemo.createdForm()
    var allData = purchasecreditmemo.dataMasterpurchasecreditmemoOriginal()
    var data = _.find(allData, function(o) {
        return o._id == e;
    });
    if (data.Currency == "IDR") {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceUSD = 0
            data.ListDetail[i].AmountUSD = 0
            purchasecreditmemo.TotalAll(data.TotalIDR)
            purchasecreditmemo.GrandTotalAll(data.GrandTotalIDR)
        }
    } else {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceIDR = 0
            data.ListDetail[i].AmountIDR = 0
            purchasecreditmemo.TotalAll(data.TotalUSD)
            purchasecreditmemo.GrandTotalAll(data.GrandTotalUSD)
        }
    }
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    ko.mapping.fromJS(data, purchasecreditmemo.record)
     if (data.DownPayment == 0){
        purchasecreditmemo.record.DownPayment("0")
     }
    if (data.VAT > 0) {
        $('#checkvat').bootstrapSwitch('state', true);
    } else {
        $('#checkvat').bootstrapSwitch('state', false);
    }
    purchasecreditmemo.record.ID(data._id)
    purchasecreditmemo.maskingMoney()
    purchasecreditmemo.reset(false)
    purchasecreditmemo.backToDraft(true)
    purchasecreditmemo.showCreate(false)
    purchasecreditmemo.deleteDraft(false)
    purchasecreditmemo.printPDF(true)
    purchasecreditmemo.text("View Draft Purchase Credit Memo")
    purchasecreditmemo.disableView()
}

purchasecreditmemo.disableView = function() {
    $(".formInput").attr("disabled", "disabled")
    $(".invhide").hide()
    // $(".btnDeleteSummary").attr("disabled", "disabled")
    // $("#buttonAdd").attr("disabled", "disabled")
    $('#datepurchase').data('kendoDatePicker').enable(false);
    var dropDown1 = $("#supliercode").data("kendoDropDownList");
    dropDown1.enable(false);
    var dropDown3 = $("#payment").data("kendoDropDownList");
    dropDown3.enable(false);
    var dropDown5 = $("#salescode").data("kendoDropDownList");
    dropDown5.enable(false);
    var dropDown4 = $("#type").data("kendoDropDownList");
    // dropDown4.enable(false);
    // $("#checkvat").bootstrapSwitch('disabled', true)
}

purchasecreditmemo.enableView = function() {
    $(".formInput").removeAttr("disabled")
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#datepurchase').data('kendoDatePicker').enable(condition);
    // $(".btnDeleteSummary").removeAttr("disabled")
    // $("#buttonAdd").removeAttr("disabled")
    $(".invhide").show()
    var dropDown1 = $("#supliercode").data("kendoDropDownList");
    dropDown1.enable(true);
    var dropDown3 = $("#payment").data("kendoDropDownList");
    dropDown3.enable(true);
    var dropDown5 = $("#salescode").data("kendoDropDownList");
    dropDown5.enable(true);
    var dropDown4 = $("#type").data("kendoDropDownList");
    // dropDown4.enable(true);
    // $("#checkvat").bootstrapSwitch('disabled', false)
}

purchasecreditmemo.resetView = function() {
    ko.mapping.fromJS(purchasecreditmemo.newRecord(), purchasecreditmemo.record)
    $(".formInput").val("")
    $("#supliername").val("")
    $(".Amount").val("")
    $("#supliercode").val(0)
    $('.nav-tabs a[href="#Create"]').tab('show')
    purchasecreditmemo.text("Create Purchase Credit Memo")
    $('#supliercode').data('kendoDropDownList').value(-1);
    $('#payment').data('kendoDropDownList').value(-1);
    // $('#type').data('kendoDropDownList').value(-1);
    purchasecreditmemo.enableView()
    // purchasecreditmemo.record.DocumentNumber(purchasecreditmemo.setDocumentNumber())
    purchasecreditmemo.maskingMoney()
    purchasecreditmemo.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    $("#checkvat").bootstrapSwitch('state', false)
    purchasecreditmemo.TotalAll(0)
    purchasecreditmemo.GrandTotalAll(0)
}

purchasecreditmemo.createdForm = function() {
    purchasecreditmemo.resetView()
    // purchasecreditmemo.record.DocumentNumber(purchasecreditmemo.setDocumentNumber())
    purchasecreditmemo.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    purchasecreditmemo.showCreate(true)
    purchasecreditmemo.reset(true)
    purchasecreditmemo.postDraft(false)
    purchasecreditmemo.printPDF(false)
    purchasecreditmemo.deleteDraft(false)
    purchasecreditmemo.disableView();
    $('#documentnumber').prop("disabled", false);
    purchasecreditmemo.backToDraft(false)
    $("#btnSave").prop("disabled", true);
}

purchasecreditmemo.editDraft = function(e) {
    purchasecreditmemo.createdForm()
    var data = _.find(purchasecreditmemo.dataMasterpurchasecreditmemo(), function(o) {
        return o._id == e;
    });
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    // data.Remark = data.Remark.substring(16)
    if (data.Currency == "IDR") {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceUSD = 0
            data.ListDetail[i].AmountUSD = 0
            purchasecreditmemo.TotalAll(data.TotalIDR)
            purchasecreditmemo.GrandTotalAll(data.GrandTotalIDR)
        }
    } else {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceIDR = 0
            data.ListDetail[i].AmountIDR = 0
            purchasecreditmemo.TotalAll(data.TotalUSD)
            purchasecreditmemo.GrandTotalAll(data.GrandTotalUSD)
        }
    }
    // data.Currency = "TEST"
    ko.mapping.fromJS(data, purchasecreditmemo.record)
    if (data.VAT > 0) {
        $('#checkvat').bootstrapSwitch('state', true);
    } else {
        $('#checkvat').bootstrapSwitch('state', false);
    }
    purchasecreditmemo.record.ID(e)
    newArr = []
    _.each(purchasecreditmemo.record.ListDetail(), function(v, i) {
        newArr.push(purchasecreditmemo.listDetail(ko.mapping.toJS(v)))
    })

    purchasecreditmemo.record.ListDetail(newArr)
    purchasecreditmemo.reset(false)
    purchasecreditmemo.showCreate(true)
    purchasecreditmemo.postDraft(true)
    purchasecreditmemo.backToDraft(true)
    purchasecreditmemo.deleteDraft(true)
    purchasecreditmemo.printPDF(false)
    $(".btnDelete").attr('disabled', false);
    $(".formInput").removeAttr("disabled")
    purchasecreditmemo.text("Edit Draft Purchase Credit Memo")
    purchasecreditmemo.enableView()
    purchasecreditmemo.checkdata()
    purchasecreditmemo.maskingMoney()
}

purchasecreditmemo.backDraft = function() {
    $('.nav-tabs a[href="#List"]').tab('show')
    $("#List").addClass("active");
    $("#Create").removeClass("active");

}
purchasecreditmemo.delete = function() {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + ko.mapping.toJS(purchasecreditmemo.record).DocumentNumber + "?",
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
            var url = "/transaction/deletedraft";
            var param = {
                id: ko.mapping.toJS(purchasecreditmemo.record).ID,
                documentNumber: ko.mapping.toJS(purchasecreditmemo.record).DocumentNumber,
            };
            ajaxPost(url, param, function(e) {
                if (e.Message == "OK") {
                    setTimeout(function() {
                        swal({
                            title: "Success!",
                            text: "Data has been deleted!",
                            type: "success",
                            confirmButtonColor: "#3da09a"
                        }, function() {
                            window.location.assign("/transaction/purchasecreditmemo")
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

purchasecreditmemo.onChangeStatus = function(textFilter) {
    purchasecreditmemo.dataMasterpurchasecreditmemo([])
    var allData = purchasecreditmemo.dataMasterpurchasecreditmemoOriginal()
    if (textFilter != "" || textFilter != undefined) {

        var Data = _.filter(allData, function(o) {
            return o.Status.indexOf(textFilter) > -1
        });
        purchasecreditmemo.dataMasterpurchasecreditmemo(Data)
    }
    purchasecreditmemo.renderGrid()

}


purchasecreditmemo.search = function(e) {
    purchasecreditmemo.textSearch(e)
    purchasecreditmemo.filterindicator(true)
    purchasecreditmemo.getDatapurchasecreditmemo(function() {
        purchasecreditmemo.renderGrid()
    })
}

purchasecreditmemo.refreshTab = function () {
    $('#dateStart').data('kendoDatePicker').value(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
    $('#dateEnd').data('kendoDatePicker').value(moment().format('YYYY-MM-DD hh:mm'))
    $("#thissupplier").data("kendoDropDownList").value(-1)
    purchasecreditmemo.filterindicator(false)
    purchasecreditmemo.getDatapurchasecreditmemo(function() {
        purchasecreditmemo.renderGrid()
    })
}

purchasecreditmemo.showData = function(e) {
    //console.log("data show data ===========" + e)
    model.Processing(true)
    var param = {}
    param = {
        TextSearch: e,
    }
    
    ko.mapping.fromJS(purchasecreditmemo.newRecord(), purchasecreditmemo.record)
    purchasecreditmemo.TotalAll(0)
    purchasecreditmemo.GrandTotalAll(0)
    purchasecreditmemo.record.VAT(0)
    purchasecreditmemo.disableView();
    $("#btnSave").prop("disabled", true);
    ajaxPost('/transaction/getdatapurchaseinvoiceinventoryforcreditmemo', param, function(res) {
        if (res.Data.length == 0) {
            swal({
                title: "Search Not Found!",
                text: "Data is not found",
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#documentnumber').val("")
            model.Processing(false)
            return
        }
        model.Processing(false)
        var data = res.Data[0]
        if (data != undefined) {
            //console.log(data)
            purchasecreditmemo.record.DatePosting(moment(data.DatePosting).format('DD-MMM-YYYY'))

            purchasecreditmemo.TotalAll(data.TotalIDR)
            purchasecreditmemo.GrandTotalAll(data.GrandTotalIDR)
            // data.Currency = "TEST"
            ko.mapping.fromJS(data, purchasecreditmemo.record)
            if (data.VAT > 0) {
                $('#checkvat').bootstrapSwitch('state', true);
            } else {
                $('#checkvat').bootstrapSwitch('state', false);
            }
            purchasecreditmemo.record.ID("")
            newArr = []
            _.each(purchasecreditmemo.record.ListDetail(), function(v, i) {
                newArr.push(purchasecreditmemo.listDetail(ko.mapping.toJS(v)))
            })

            purchasecreditmemo.record.ListDetail(newArr)
            purchasecreditmemo.reset(false)
            purchasecreditmemo.showCreate(true)
            purchasecreditmemo.postDraft(true)
            purchasecreditmemo.backToDraft(true)
            purchasecreditmemo.deleteDraft(true)
            purchasecreditmemo.printPDF(false)
            $(".btnDelete").attr('disabled', false);
            // console.log("disable")
            // $(".formInput").removeAttr("disabled")
            // disable view
            $(".formInput").attr("disabled", "disabled")
            $(".priceidr").attr("disabled", "disabled")
            $("[name='downPayment']").removeAttr("disabled")
            $("[name='remark']").removeAttr("disabled")
            $("#checkvat").bootstrapSwitch('disabled', true)
            var dropDown1 = $("#supliercode").data("kendoDropDownList");
            dropDown1.enable(false);
            var dropDown3 = $("#payment").data("kendoDropDownList");
            dropDown3.enable(false);
            $("#buttonAdd").attr("disabled", "disabled")
            if (purchasecreditmemo.record.VAT() > 0) {
                $("[name='pay']").bootstrapSwitch({
                    disabled: true,
                    state: true
                });
            }
            $("#btnSave").prop("disabled", false);
            $("#btnReset").show();
            $("#btnDelete").hide();
            purchasecreditmemo.disableView();
            purchasecreditmemo.record.Remark(data.DocumentNumber + " - ")
            $("#Remark").prop("disabled", false).focus();

            purchasecreditmemo.text("Edit Draft Purchase Credit Memo")
            // purchasecreditmemo.enableView()
            // purchasecreditmemo.checkdata()
            purchasecreditmemo.maskingMoney()
        }

    })
}

purchasecreditmemo.init = function() {
    ProActive.KendoDatePickerRange();
    purchasecreditmemo.getDatapurchasecreditmemo(function() {
        purchasecreditmemo.renderGrid()
    })
    purchasecreditmemo.maskingMoney()
    purchasecreditmemo.getDateNow()
    purchasecreditmemo.switchButton()
    purchasecreditmemo.getDataSales()
    purchasecreditmemo.setDate()
    purchasecreditmemo.GetDataPurchaseInvoice()
}
purchasecreditmemo.fromPOINVSummary = function() {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    if (num != null) {
        var allData = purchasecreditmemo.dataMasterpurchasecreditmemoOriginal()
        var data = _.find(allData, function(o) {
            return o.DocumentNumber == num;
        });
        if (data != undefined) {
            purchasecreditmemo.viewDraft(data._id)
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning",
                confirmButtonColor: "#3da09a"
            }, function() {
                window.location.assign("/transaction/purchasecreditmemo")
            });
        }
    }
}

purchasecreditmemo.filterText = function (term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["DocumentNumber", "SupplierCode", "SupplierName", "Remark"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator: "contains", value: $searchValue });
    }
    $("#gridListPurchasecreditmemo").data("kendoGrid").dataSource.query({ filter: filter });
}

$(function() {

    // purchasecreditmemo.getDocNumberPO()
    purchasecreditmemo.getDataSupplier()
    purchasecreditmemo.getDataAccount()
    purchasecreditmemo.getDataTypePurchase()
    purchasecreditmemo.init()
    purchasecreditmemo.text("Created Purchase Credit Memo")
    $(".priceusd").keyup(function() {
        var nameLength = $(".priceusd").val().length;
        if (nameLength > 0) {
            $(".priceidr").attr("disabled", "disabled");
        } else {
            $(".priceidr").removeAttr("disabled");
        }
    });
    $(".priceidr").keyup(function() {
        var nameLength = $(".priceidr").val().length;
        if (nameLength > 0) {
            $(".priceusd").attr("disabled", "disabled");
        } else {
            $(".priceusd").removeAttr("disabled");
        }
    });

    $("#textSearch").on("keyup blur change", function () {
        purchasecreditmemo.filterText();
    });

    $('#documentnumber').keydown(function(e) {
        var key = e.charCode ? e.charCode : e.keyCode ? e.keyCode : 0;
        if (key == 13) {
            e.preventDefault();
            //console.log($('#documentnumber').val())
            purchasecreditmemo.showData($('#documentnumber').val());
        }
    });
    purchasecreditmemo.initAutoComplete();
}); 

purchasecreditmemo.retrievedAutoCompleteData = [];
purchasecreditmemo.initAutoComplete = function() {
    purchasecreditmemo.autoCompleteDataSource = new kendo.data.DataSource({
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
                var url = "/transaction/getautocinvnumpcm";
                ajaxPost(url, e.data.filter, function (res) {
                    purchasecreditmemo.retrievedAutoCompleteData = res.data;
                    e.success(res);
                }, function () {
                    e.error();
                })
            }
        }
    });
    $("#documentnumber").kendoAutoComplete({
        dataSource: purchasecreditmemo.autoCompleteDataSource,
        filter: "contains",
        clearButton: false,
        placeholder: "Enter Invoice Number...",
        minLength: 12,
        dataValueField : "DocumentNumber",
        dataTextField : "DocumentNumber",
    });
}