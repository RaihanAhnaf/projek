model.Processing(false)

var purchaseorder = {}
purchaseorder.dataMasterSupplier = ko.observableArray([])
purchaseorder.dataDropDownSupplier = ko.observableArray([])
purchaseorder.dataDropDownSupplierFilter = ko.observableArray([])
purchaseorder.dataMasterAccount = ko.observableArray([])
purchaseorder.dataDropDownAccount = ko.observableArray([])
purchaseorder.dataMasterPurchaseOrder = ko.observableArray([])
purchaseorder.dataListPurchaseOrder = ko.observableArray([])
purchaseorder.dataListDraftPurchaseOrder = ko.observableArray([])
purchaseorder.dataMasterPurchaseOrderOriginal = ko.observableArray([])
purchaseorder.type = ko.observableArray([])
purchaseorder.checkedvat = ko.observable(false)
purchaseorder.text = ko.observable()
purchaseorder.textSearch = ko.observable()
purchaseorder.statustext = ko.observable()
purchaseorder.sequenceNumber = ko.observable()
purchaseorder.showCreate = ko.observable(true)
purchaseorder.backToDraft = ko.observable(false)
purchaseorder.postDraft = ko.observable(false)
purchaseorder.deleteDraft = ko.observable(false)
purchaseorder.printPDF = ko.observable(false)
purchaseorder.reset = ko.observable(false)
purchaseorder.setDocumentNumber = ko.observable('')
purchaseorder.rowSpan = ko.observable(5)
purchaseorder.DatePageBar = ko.observable()
purchaseorder.BoolVat = ko.observable(false)
purchaseorder.textSupplierSearch = ko.observable()
purchaseorder.filterindicator = ko.observable(false)
purchaseorder.filterindicatorInv = ko.observable(false)
purchaseorder.TotalAll = ko.observable(0)
purchaseorder.GrandTotalAll = ko.observable(0)
purchaseorder.names = ko.observable('')
purchaseorder.dataMasterSales = ko.observableArray([])
purchaseorder.dataDropDownSales = ko.observableArray([])
purchaseorder.dataMasterInventory = ko.observableArray([])
purchaseorder.dataDropDownInventory = ko.observableArray([])
purchaseorder.dataMasterPurchaseInventory = ko.observableArray([])
purchaseorder.dataMasterPurchaseInventoryOriginal = ko.observableArray([])
purchaseorder.dataSupplierAll = [];
purchaseorder.filterAlert = ko.observable(false) 
purchaseorder.filterAlertType = ko.observable("search")
purchaseorder.payment = [{
    value: "CASH",
    text: "Cash"
}, {
    value: "INSTALLMENT",
    text: "Installment"
}]

purchaseorder.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    purchaseorder.DatePageBar(page)
}
purchaseorder.newRecord = function () {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumber: "",
        SupplierCode: "",
        SupplierName: "",
        AccountCode: 0,
        Payment: "",
        Type: "",
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
        SalesCode: "",
        SalesName: "",
        LocationID: 0,
        LocationName: ""
    }
    page.ListDetail.push(purchaseorder.listDetail({}))
    return page
}

purchaseorder.listDetail = function (data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.Item = data.Item == undefined ? '' : data.Item
    dataTmp.Qty = data.Qty == undefined ? '' : data.Qty
    dataTmp.PriceUSD = data.PriceUSD == undefined ? '' : data.PriceUSD
    dataTmp.PriceIDR = data.PriceIDR == undefined ? '' : data.PriceIDR
    dataTmp.AmountUSD = data.AmountUSD == undefined ? '' : data.AmountUSD
    dataTmp.AmountIDR = data.AmountIDR == undefined ? '' : data.AmountIDR
    if (purchaseorder.names() == "Inventory") {
        dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
    }

    var x = ko.mapping.fromJS(dataTmp)

    x.Qty.subscribe(function (e) {
        //Menghitung Total Amount IDR
        total = FormatCurrency(e) * FormatCurrency(x.PriceIDR())
        x.AmountIDR(total)

        //Menghitung Total Amount USD
        totalUSD = FormatCurrency(e) * FormatCurrency(x.PriceUSD())
        x.AmountUSD(totalUSD)


        //Menghitung Total
        var data = ko.mapping.toJS(purchaseorder.record)
        if (purchaseorder.record.ListDetail()[0].PriceIDR > 0 || total > 0) {
            var alldataidr = ko.mapping.toJS(purchaseorder.record.ListDetail())
            var totalIDR = _.sumBy(alldataidr, function (v) {
                return v.AmountIDR
            })
            purchaseorder.record.TotalIDR(totalIDR)
            purchaseorder.TotalAll(ChangeToRupiah(purchaseorder.record.TotalIDR()))

        } else {
            var alldatausd = ko.mapping.toJS(purchaseorder.record.ListDetail())
            var totalUSD = _.sumBy(alldatausd, function (v) {
                return v.AmountUSD
            })
            purchaseorder.record.TotalUSD(totalUSD)
            var totalIDRRATE = totalUSD * FormatCurrency(purchaseorder.record.Rate())
            purchaseorder.record.TotalIDR(totalIDRRATE)
            purchaseorder.TotalAll(ChangeToRupiah(purchaseorder.record.TotalUSD()))

        }

    })

    x.PriceIDR.subscribe(function (e) {
        x.AmountIDR(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var alldataidr = ko.mapping.toJS(purchaseorder.record.ListDetail())
        var totalIDR = _.sumBy(alldataidr, function (v) {
            return v.AmountIDR
        })
        // purchaseorder.record.Total(ChangeToRupiah(totalIDR))
        purchaseorder.record.TotalIDR(totalIDR)
        purchaseorder.record.Currency("IDR")
        purchaseorder.TotalAll(ChangeToRupiah(purchaseorder.record.TotalIDR()))
        purchaseorder.checkdata()

    })

    x.PriceUSD.subscribe(function (e) {
        x.AmountUSD(FormatCurrency(e) * FormatCurrency(x.Qty()))
        var data = ko.mapping.toJS(purchaseorder.record)
        // x.PriceIDR(FormatCurrency(e)*FormatCurrency(data.Rate))
        var alldatausd = ko.mapping.toJS(purchaseorder.record.ListDetail())
        var totalUSD = _.sumBy(alldatausd, function (v) {
            return v.AmountUSD
        })
        // purchaseorder.record.Total(ChangeToRupiah(totalIDR))
        purchaseorder.record.TotalUSD(totalUSD)
        var totalIDRRATE = totalUSD * FormatCurrency(purchaseorder.record.Rate())
        purchaseorder.record.TotalIDR(totalIDRRATE)
        purchaseorder.record.Currency("USD")
        purchaseorder.TotalAll(ChangeToRupiah(purchaseorder.record.TotalUSD()))
        purchaseorder.checkdata()
    })

    return x
}

purchaseorder.record = ko.mapping.fromJS(purchaseorder.newRecord())

purchaseorder.switchButton = function () {
    $('#checkvat').on('switchChange.bootstrapSwitch', function (event, state) {
        var data = ko.mapping.toJS(purchaseorder.record)
        purchaseorder.BoolVat(state)
        if (state) {
            if (purchaseorder.record.Currency() == "IDR") {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalIDR)
                purchaseorder.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
                var Totalall = GT + totVat
                purchaseorder.record.GrandTotalIDR(Totalall)
                purchaseorder.GrandTotalAll(ChangeToRupiah(Totalall))

            } else {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalUSD)
                purchaseorder.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
                var Totalall = GT + totVat
                purchaseorder.record.GrandTotalUSD(Totalall)
                purchaseorder.record.GrandTotalIDR(Totalall * FormatCurrency(data.Rate))
                purchaseorder.GrandTotalAll(ChangeToRupiah(Totalall))

            }
        } else {
            purchaseorder.record.VAT(ChangeToRupiah(0))
            if (purchaseorder.record.Currency() == "IDR") {
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
                purchaseorder.record.GrandTotalIDR(GT)
                purchaseorder.GrandTotalAll(ChangeToRupiah(GT))

            } else {
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
                purchaseorder.record.GrandTotalUSD(GT)
                purchaseorder.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
                purchaseorder.GrandTotalAll(ChangeToRupiah(GT))

            }
        }
    });
}

purchaseorder.TotalAll.subscribe(function (e) {
    var data = ko.mapping.toJS(purchaseorder.record)
    if (purchaseorder.BoolVat()) {
        var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(e)
        var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(e)
        purchaseorder.record.VAT(totVat)
        if (purchaseorder.record.Currency() == "IDR") {
            purchaseorder.record.GrandTotalIDR(GT + totVat)
            purchaseorder.GrandTotalAll(ChangeToRupiah(GT + totVat))
        } else {
            purchaseorder.record.GrandTotalUSD(GT + totVat)
            purchaseorder.record.GrandTotalIDR((GT + totVat) * FormatCurrency(data.Rate))
            purchaseorder.GrandTotalAll(ChangeToRupiah(GT + totVat))
        }
    } else {
        purchaseorder.record.VAT(0)
        if (data.ListDetail.length != 0) {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(e)
            if (purchaseorder.record.Currency() == "IDR") {
                purchaseorder.record.GrandTotalIDR(GT)
                purchaseorder.GrandTotalAll(ChangeToRupiah(GT))
            } else {
                purchaseorder.record.GrandTotalUSD(GT)
                purchaseorder.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
                purchaseorder.GrandTotalAll(ChangeToRupiah(GT))
            }

        } else {
            purchaseorder.GrandTotalAll(0)
        }
    }
})

purchaseorder.record.Discount.subscribe(function (e) {
    var data = ko.mapping.toJS(purchaseorder.record)
    if (FormatCurrency(e) > 100 || FormatCurrency(e) < 0) {
        return swal({
            title: "Warning!",
            text: "your discount are irational",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function () {
            purchaseorder.record.Discount(0)
        });
    }


    if (purchaseorder.record.Currency() == "IDR") {
        var GT = ((100 - FormatCurrency(e)) / 100) * FormatCurrency(data.TotalIDR)
        var totVat = 0
        if (purchaseorder.BoolVat()) {
            totVat = ((100 - FormatCurrency(e)) / 1000) * FormatCurrency(data.TotalIDR)
            purchaseorder.record.VAT(ChangeToRupiah(totVat))
        } else {
            purchaseorder.record.VAT(ChangeToRupiah(0))
            totVat = 0
        }
        purchaseorder.record.GrandTotalIDR(GT + totVat)
        purchaseorder.GrandTotalAll(ChangeToRupiah(GT + totVat))
    } else {
        var GT = ((100 - FormatCurrency(e)) / 100) * FormatCurrency(data.TotalUSD)
        var totVat = 0
        if (purchaseorder.BoolVat()) {
            var totVat = ((100 - FormatCurrency(e)) / 1000) * FormatCurrency(data.TotalUSD)
            purchaseorder.record.VAT(ChangeToRupiah(totVat))
        } else {
            purchaseorder.record.VAT(ChangeToRupiah(0))
            totVat = 0
        }

        purchaseorder.record.GrandTotalUSD(GT + totVat)
        purchaseorder.record.GrandTotalIDR((GT + totVat) * FormatCurrency(data.Rate))
        purchaseorder.GrandTotalAll(ChangeToRupiah(GT + totVat))

    }

})

purchaseorder.record.VAT.subscribe(function (e) {
    var data = ko.mapping.toJS(purchaseorder.record)
    // purchaseorder.BoolVat(state)
    if (purchaseorder.BoolVat()) {
        if (purchaseorder.record.Currency() == "IDR") {
            var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalIDR)
            purchaseorder.record.VAT(ChangeToRupiah(totVat))
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
            var Totalall = GT + totVat
            purchaseorder.record.GrandTotalIDR(Totalall)
            purchaseorder.GrandTotalAll(ChangeToRupiah(Totalall))

        } else {
            var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalUSD)
            purchaseorder.record.VAT(ChangeToRupiah(totVat))
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
            var Totalall = GT + totVat
            purchaseorder.record.GrandTotalUSD(Totalall)
            purchaseorder.record.GrandTotalIDR(Totalall * FormatCurrency(data.Rate))
            purchaseorder.GrandTotalAll(ChangeToRupiah(Totalall))

        }
    } else {
        purchaseorder.record.VAT(ChangeToRupiah(0))
        if (purchaseorder.record.Currency() == "IDR") {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
            purchaseorder.record.GrandTotalUSD(GT)
            purchaseorder.GrandTotalAll(ChangeToRupiah(GT))

        } else {
            var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
            purchaseorder.record.GrandTotalUSD(GT)
            purchaseorder.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
            purchaseorder.GrandTotalAll(ChangeToRupiah(GT))

        }
    }
})
purchaseorder.record.DownPayment.subscribe(function (e) {
    if (parseFloat(e) > 100 || parseFloat(e) < 0) {
        return swal({
            title: "Warning!",
            text: "your Down Payment are irational",
            type: "info",
            confirmButtonColor: "#3da09a"
        }, function () {
            purchaseorder.record.DownPayment(0)
        });
    }
})

// purchaseorder.record.Rate.subscribe(function (e) {
//     var data = ko.mapping.toJS(purchaseorder.record)
//     if (data.ListDetail[0].PriceUSD > 0) {
//         for (var i = 0; i < data.ListDetail.length; i++) {
//             var alldatausd = ko.mapping.toJS(purchaseorder.record.ListDetail())
//             var totalIDR = _.sumBy(alldatausd, function (v) {
//                 return v.AmountUSD * FormatCurrency(e)
//             })
//             purchaseorder.record.Total(ChangeToRupiah(totalIDR))
//         }
//     }
// })
purchaseorder.getBalance = function(code, type, callback) {
    var p = {
        CustomerCode : code
    }

    ajaxPost((type == "CUSTOMER" ? '/master/GetCustomerBalance' : '/master/GetSupplierBalance').toLowerCase(), p, function (res) {
        if (typeof callback == "function") {
            callback(res.Data);
        }
    })
}

purchaseorder.getDataSupplier = function () {
    model.Processing(true)

    ajaxPost('/transaction/getsupplier', {}, function (res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseorder.dataMasterSupplier(res.Data)
        purchaseorder.dataDropDownSupplier(res.Data)
        var DataSupplier = res.Data
        purchaseorder.dataSupplierAll = res.Data;
        for (i in DataSupplier) {
            purchaseorder.getBalance(DataSupplier[i].Kode, "SUPPLIER", function(data){
                this.Balance = data.BalanceIDR;
            }.bind(DataSupplier[i]));
            DataSupplier[i].Kode = DataSupplier[i].Kode
            DataSupplier[i].KodeName = DataSupplier[i].Kode + " - " + DataSupplier[i].Name
            DataSupplier[i].Name = DataSupplier[i].Name
        }
        purchaseorder.dataDropDownSupplierFilter(DataSupplier)

        model.Processing(false)
    })
}

purchaseorder.getDataTypePurchase = function () {
    model.Processing(true)

    ajaxPost('/transaction/gettypepurchase', {}, function (res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseorder.type(res.Data)
        model.Processing(false)
    })
}

purchaseorder.getDataAccount = function () {
    model.Processing(true)

    ajaxPost('/transaction/getaccount', {}, function (res) {

        if (res.Total === 0) {
            // console.log(res)
            swal({
                title: "Info!",
                text: "Chart of account is empty, Please input chart of account!",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseorder.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        purchaseorder.dataDropDownAccount(DataAccount)

        model.Processing(false)
    })
}

purchaseorder.onChangeAccountNumber = function (value) {
    var result = _.filter(purchaseorder.dataMasterSupplier(), {
        'Kode': value
    })[0].Name
    purchaseorder.record.SupplierName(result);
}
purchaseorder.onChangeSalesCode = function (value) {
    if (value != "") {
        var result = _.find(purchaseorder.dataMasterSales(), {
            'Kode': value
        }).SalesName
        purchaseorder.record.SalesName(result);
    }
}
purchaseorder.maskingMoney = function () {
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

purchaseorder.getDataPurchaseOrder = function (callback) {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {}
    if (purchaseorder.filterindicator() == true) {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true,
            TextSearch: purchaseorder.textSearch(),
            SupplierCode: purchaseorder.textSupplierSearch(),
        }
    } else {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: false,
        }
    }

    model.Processing(true)
    ajaxPost('/transaction/getdatapurchaseorderstateless', param, function (res) {
        if (res.IsError && purchaseorder.filterAlert()) {
            swal({
                title: "No Data Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            res.Data = []
        }
        purchaseorder.dataMasterPurchaseOrder(res.Data)
        purchaseorder.dataMasterPurchaseOrderOriginal(res.Data)
        model.Processing(false)
        callback()
        purchaseorder.fromPOINVSummary()
    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}
purchaseorder.getDataPurchaseInventory = function (callback) {
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();

    var param = {}
    if (purchaseorder.filterindicatorInv() == true) {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: true,
            TextSearch: purchaseorder.textSearch(),
            SupplierCode: purchaseorder.textSupplierSearch(),
        }
    } else {
        param = {
            DateStart: moment(startdate).format('YYYY-MM-DD'),
            DateEnd: moment(enddate).format('YYYY-MM-DD'),
            Filter: false,
        }
    }

    ajaxPost('/transaction/getdatapurchaseinventorystateless', param, function (res) {
        if (res.IsError&& purchaseorder.filterAlert()) {
            swal({
                title: "No Data Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            res.Data = []
            //  return
        }
        purchaseorder.dataMasterPurchaseInventory(res.Data)
        purchaseorder.dataMasterPurchaseInventoryOriginal(res.Data)
        model.Processing(false)
        callback()
        purchaseorder.fromPOINVSummary()
    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

purchaseorder.renderGrid = function () {
    var data = purchaseorder.dataMasterPurchaseOrder();
    // for (var i = 0; i < data.le.Kondisingth; i++) {
    //     data[i].Kondisi = (data[i].User ==  )
    // }
    /*
    if (typeof $('#gridListPurchaseOrder').data('kendoGrid') !== 'undefined') {
        $('#gridListPurchaseOrder').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
        }))
        return
    }
    */

    var columns = [

        {
            title: 'Action',
            width: 100,
            template: "# if ((User == userinfo.usernameh() || userinfo.usernameh() == 'administrator' || userinfo.rolenameh() == 'supervisor') && Status == 'PO' ) {#<button onclick='purchaseorder.viewDraft(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button onclick='purchaseorder.editDraft(\"#: ID #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button>#}else{#<button onclick='purchaseorder.viewDraft(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button class='btn btn-sm btn-success btn-flat pointer-disable' style='display:none'><i class='fa fa-pencil'></i></button>#}#",
        }, {
            field: 'DateStr',
            title: 'Date',
            width: 120,
        }, {
            field: 'DocumentNumber',
            title: 'Order #',
            width: 160,
        }, {
            field: 'SupplierName',
            title: 'Supplier Name',
            width: 150,
        }, {
            title: 'Location',
            width: 120,
            template: "#: LocationID # - #: LocationName #",
        }, {
            field: '',
            title: 'Order Total ',
            width: 180,
            // template: "#=ChangeToRupiah(GrandTotalUSD)#"
            attributes: {
                "class": "rightAction",
            },
            template: "#if ( Currency == 'USD') {# $ #: ChangeToRupiah(GrandTotalUSD) # #} else {# Rp. #: ChangeToRupiah(GrandTotalIDR) # #}#"

        }, {
            field: 'Status',
            title: 'Status',
            width: 50,
        }, {
            field: 'Remark',
            title: 'Remark',
            width: 200,
        }
    ]

    $('#gridListPurchaseOrder').kendoGrid({
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
        excelExport: function (e) {
            ProActive.kendoExcelRender(e, "PurchaseOrderNonInventory", function (row, sheet) {
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "data") {
                    }
                }
            });
        },
    })
}

purchaseorder.renderGridInventory = function () {
    var data = purchaseorder.dataMasterPurchaseInventory();
    // for (var i = 0; i < data.le.Kondisingth; i++) {
    //     data[i].Kondisi = (data[i].User ==  )
    // }
    /*
    if (typeof $('#gridListPurchaseOrder').data('kendoGrid') !== 'undefined') {
        $('#gridListPurchaseOrder').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
        }))
        return
    }
    */

    var columns = [

        {
            title: 'Action',
            width: 100,
            template: "# if ((User == userinfo.usernameh() || userinfo.usernameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' ) && Status == 'PO') {#<button onclick='purchaseorder.viewDraft(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button onclick='purchaseorder.editDraft(\"#: ID #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button>#}else{#<button onclick='purchaseorder.viewDraft(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> <button class='btn btn-sm btn-success btn-flat pointer-disable' style='display:none'><i class='fa fa-pencil'></i></button>#}#",
        }, {
            field: 'DateStr',
            title: 'Date',
            width: 120,
        }, {
            field: 'DocumentNumber',
            title: 'Order #',
            width: 160,
        }, {
            field: 'SupplierName',
            title: 'Supplier Name',
            width: 150,
        }, {
            title: 'Location',
            width: 120,
            template: "#: LocationID # - #: LocationName #",
        }, {
            field: '',
            title: 'Order Total ',
            width: 180,
            // template: "#=ChangeToRupiah(GrandTotalUSD)#"
            attributes: {
                "class": "rightAction",
            },
            template: "#if ( Currency == 'USD') {# $ #: ChangeToRupiah(GrandTotalUSD) # #} else {# Rp. #: ChangeToRupiah(GrandTotalIDR) # #}#"

        }, {
            field: 'Status',
            title: 'Status',
            width: 50,
        }, {
            field: 'Remark',
            title: 'Remark',
            width: 200,
        }
    ]

    $('#gridListPurchaseOrder').kendoGrid({
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
        excelExport: function (e) {
            ProActive.kendoExcelRender(e, "PurchaseOrderInventory", function (row, sheet) {
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "data") {
                    }
                }
            });
        },
    })
}
purchaseorder.exportExcel = function () {
    $("#gridListPurchaseOrder").getKendoGrid().saveAsExcel();
}


purchaseorder.checkdata = function () {
    if (FormatCurrency(purchaseorder.record.ListDetail()[0].PriceUSD()) == 0 && FormatCurrency(purchaseorder.record.ListDetail()[0].PriceIDR()) == 0) {
        $(".priceidr").removeAttr("disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(purchaseorder.record.ListDetail()[0].PriceUSD()) > 0) {
        $(".priceidr").attr("disabled", "disabled");
        $(".priceusd").removeAttr("disabled");
    } else if (FormatCurrency(purchaseorder.record.ListDetail()[0].PriceUSD()) == 0) {
        $(".priceusd").attr("disabled", "disabled");
        $(".priceidr").removeAttr("disabled");
    }

}

purchaseorder.addNewItem = function () {
    purchaseorder.record.ListDetail.push(purchaseorder.listDetail({}))
    purchaseorder.checkdata()
    purchaseorder.maskingMoney()
    if (purchaseorder.names() == "Non Inventory") {
        $(".invhide").hide()
        $(".item").prop('disabled', false);

    } else {
        $(".invhide").show()
        $(".item").prop('disabled', true);
    }
}

purchaseorder.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}


purchaseorder.removeRow = function () {
    purchaseorder.checkdata()
    purchaseorder.record.ListDetail.remove(this)
    if (purchaseorder.record.ListDetail().length == 0) {
        purchaseorder.record.ListDetail.push(purchaseorder.listDetail({}))
    }
    //Menghitung Total 
    var alldataidr = ko.mapping.toJS(purchaseorder.record.ListDetail())
    var totalIDR = _.sumBy(alldataidr, function (v) {
        return v.AmountIDR
    })

    var alldata = ko.mapping.toJS(purchaseorder.record.ListDetail())
    var totalUSD = _.sumBy(alldata, function (v) {
        return v.AmountUSD
    })
    if (FormatCurrency(purchaseorder.record.ListDetail()[0].PriceIDR()) > 0) {
        purchaseorder.record.TotalIDR(ChangeToRupiah(totalIDR))
    }
    //  if (FormatCurrency(purchaseorder.record.ListDetail()[0].PriceUSD()) > 0) {
    //      purchaseorder.record.Total(ChangeToRupiah(totalUSD))
    //  }

    var data = ko.mapping.toJS(purchaseorder.record)
    //Menghitung VAT           
    var nomvat = ((100 - data.Discount) / 1000) * FormatCurrency(data.Total)
    purchaseorder.record.VAT(ChangeToRupiah(nomvat))

    // Menghitung GrandTotal
    var totalnumber = FormatCurrency(data.Total)
    var nomgrandtotal = (((100 - data.Discount) / 100) * totalnumber) + nomvat
    purchaseorder.record.GrandTotalIDR(ChangeToRupiah(nomgrandtotal))
    if (purchaseorder.names() == "Non Inventory") {
        $(".invhide").hide()
        $(".item").prop('disabled', false);

    } else {
        $(".invhide").show()
        $(".item").prop('disabled', true);
    }

}
purchaseorder.getDocNumberPO = function () {
    model.Processing(true)
    ajaxPost('/transaction/getlastnumberpo', {}, function (res) {
        var tgl = moment(Date()).format('/DDMMYY/')
        var zr
        if (res.Number < 10) {
            zr = '00'
        } else if (res.Number >= 10 && res.Number < 100) {
            zr = '0'
        } else {
            zr = ''
        }

        purchaseorder.record.DocumentNumber('PO' + tgl + zr + res.Number);
        purchaseorder.sequenceNumber(res.Number)
        purchaseorder.setDocumentNumber('PO' + tgl + zr + res.Number);
    })
}

purchaseorder.saveSwitch = function () {
    var validator = $("#Create").kendoValidator().data("kendoValidator")
    // console.log(validator.validate())
    if (validator.validate()) {
        if (purchaseorder.names() == "Non Inventory") {
            purchaseorder.saveData()
        } else {
            purchaseorder.saveDataInventory()
        }
    }
}

purchaseorder.saveData = function () {

    var change = ko.mapping.toJS(purchaseorder.record)
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

    if (change.Type == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't choose the type",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }

    if (change.DownPayment == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't input the Down Payment",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    // if (FormatCurrency(purchaseorder.TotalAll) == 0) {
    //     return swal('Warning!', "You didn't buy anything", "info")
    // }
    if (change.Remark == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't input the remark",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
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
    }
    change.DatePosting = moment($('#datepurchase').data('kendoDatePicker').value()).startOf("day").format("YYYY-MM-DD") + "T00:00:00Z"; 
    change.DateStr = moment($('#datepurchase').data('kendoDatePicker').value()).startOf("day").format("DD-MMM-YYYY");
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
    change.AccountCode = 0
    var param = {
        Data: change,
        // LastNumber: purchaseorder.sequenceNumber(),
    }
    var url = "/transaction/insertdraftpurchaseorder"
    // Shorthand for swal confirm submit
    ProActive.swalConfirmSubmit("Purchase Order Non Inventory", function () {
        model.Processing(true)
        ajaxPost(url, param, function (e) {
            setTimeout(function () {
                swal({
                    title: "Success!",
                    text: "Data has been saved!",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                }, function () {
                    purchaseorder.resetView()
                    purchaseorder.getDataPurchaseOrder(function () {
                        purchaseorder.renderGrid()
                    })
                    // window.location.assign("/transaction/purchaseorder")
                    // $("#btnSave").attr("disabled", "disabled")
                });
            }, 100)
            model.Processing(false)
        });
    });
    /*
    swal({
        title: "Are you sure?",
        text: "You will submit this Purchase Order",
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
                        window.location.assign("/transaction/purchaseorder")
                        $("#btnSave").attr("disabled", "disabled")
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
    */
}

purchaseorder.saveDataInventory = function () {
    var change = ko.mapping.toJS(purchaseorder.record)
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
    if (change.DownPayment == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't input the Down Payment",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (change.Remark == "") {
        return swal({
            title: 'Warning!',
            text: "You haven't input the remark",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
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
    }
    change.DatePosting = moment($('#datepurchase').data('kendoDatePicker').value()).startOf("day").format("YYYY-MM-DD") + "T00:00:00Z";
    change.DateStr = moment($('#datepurchase').data('kendoDatePicker').value()).startOf("day").format("DD-MMM-YYYY");
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
    change.AccountCode = 1401
    var param = {
        Data: change,
        // LastNumber: purchaseorder.sequenceNumber(),
    }
    var url = "/transaction/insertdraftpurchaseorderinventory"
    swal({
        title: "Are you sure?",
        text: "You will submit this Purchase Order Inventory",
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
                        // window.location.assign("/transaction/purchaseorder")
                        purchaseorder.resetView()
                        // purchaseorder.choose("Inventory")
                        // purchaseorder.backDraft()                       
                        // $("#btnSave").attr("disabled", "disabled")
                        purchaseorder.getDataPurchaseInventory(function () {
                            purchaseorder.renderGridInventory()
                        })
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


purchaseorder.printToPdf = function () {
    model.Processing(true)
    var param = {
        Id: purchaseorder.record.ID(),
    }
    var url = "/transaction/exporttopdfpurchaseorder"
    if (purchaseorder.names() == "Inventory") {
        url = "/transaction/exporttopdfpurchaseorderinventory"
    }
    ajaxPost(url, param, function (e) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/purchaseorder/' + e, '_blank');
        tabOrWindow.focus();
        // console.log(e);
        // window.open(e);
        // var pom = document.createElement('a');
        // pom.setAttribute('href', "/res/docs/purchaseorder/" + e);
        // pom.setAttribute('download', e);
        // pom.click();

        // console.log(pom)
        /*var url = ("{{BaseUrl}}/res/docs/purchaseorder/" + e)
        var link = document.createElement('a');
        link.href = url;
        link.dispatchEvent(new MouseEvent('click'));*/
        /*var fileURL = '{{BaseUrl}}res/docs/purchaseorder/' + e
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

purchaseorder.viewDraft = function (e) {
    //console.log(purchaseorder.names())
    purchaseorder.createdForm()
    var data = []
    if (purchaseorder.names() == "Non Inventory") {
        setTimeout(function () { $(".invhide").hide(); }, 50);
        // $(".invhide").hide()
        $(".item").prop('disabled', false);
        data = _.find(purchaseorder.dataMasterPurchaseOrder(), function (o) {
            return o.ID == e;
        });
        //    console.log(data)

    } else {
        $(".invhide").show()
        $(".item").prop('disabled', true);
        data = _.find(purchaseorder.dataMasterPurchaseInventory(), function (o) {
            return o.ID == e;
        });
    }
    //  console.log(data)

    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    // data.Remark = data.Remark.substring(16)
    if (data.Currency == "IDR") {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceUSD = 0
            data.ListDetail[i].AmountUSD = 0
            purchaseorder.TotalAll(data.TotalIDR)
            purchaseorder.GrandTotalAll(data.GrandTotalIDR)
        }
    } else {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceIDR = 0
            data.ListDetail[i].AmountIDR = 0
            purchaseorder.TotalAll(data.TotalUSD)
            purchaseorder.GrandTotalAll(data.GrandTotalUSD)
        }
    }
    // data.Currency = "TEST"
    ko.mapping.fromJS(data, purchaseorder.record)
    if (data.DownPayment == 0) {
        purchaseorder.record.DownPayment("0")
    }
    if (data.VAT > 0) {
        $('#checkvat').bootstrapSwitch('state', true);
    } else {
        $('#checkvat').bootstrapSwitch('state', false);
    }
    purchaseorder.record.ID(e)
    newArr = []
    _.each(purchaseorder.record.ListDetail(), function (v, i) {
        if (v != undefined) {
            newArr.push(purchaseorder.listDetail(ko.mapping.toJS(v)))
        }
    })

    purchaseorder.record.ListDetail(newArr)
    purchaseorder.record.ID(data.ID)
    purchaseorder.maskingMoney()
    purchaseorder.reset(false)
    purchaseorder.backToDraft(true)
    purchaseorder.showCreate(false)
    purchaseorder.deleteDraft(false)
    purchaseorder.printPDF(true)
    purchaseorder.text("View Draft Purchase Order")
    purchaseorder.disableView()
}

purchaseorder.disableView = function () {
    $(".formInput").attr("disabled", "disabled")
    // $(".btnDeleteSummary").attr("disabled", "disabled")
    $("#buttonAdd").hide()
    $(".btnDelete").hide()
    $(".btnhide").hide()
    $('#datepurchase').data('kendoDatePicker').enable(false);
    var dropDown1 = $("#supliercode").data("kendoDropDownList");
    dropDown1.enable(false);
    var dropDown2 = $("#salescode").data("kendoDropDownList");
    dropDown2.enable(false);
    var dropDown3 = $("#payment").data("kendoDropDownList");
    dropDown3.enable(false);
    var dropDown4 = $("#typepurchase").data("kendoDropDownList");
    dropDown4.enable(false);
    $("#checkvat").bootstrapSwitch('disabled', true)

    if (purchaseorder.names() != "Non Inventory") {
        $(".codeitem").off('click');
    }

}

purchaseorder.enableView = function () {
    $(".formInput").removeAttr("disabled")
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#datepurchase').data('kendoDatePicker').enable(condition);
    // $(".btnDeleteSummary").removeAttr("disabled")
    // $("#buttonAdd").removeAttr("disabled")
    $(".btnhide").show()
    $("#buttonAdd").show()
    $(".btnDelete").show()
    var dropDown1 = $("#supliercode").data("kendoDropDownList");
    dropDown1.enable(true);
    var dropDown2 = $("#salescode").data("kendoDropDownList");
    dropDown2.enable(true);
    var dropDown3 = $("#payment").data("kendoDropDownList");
    dropDown3.enable(true);
    var dropDown4 = $("#typepurchase").data("kendoDropDownList");
    dropDown4.enable(true);
    $("#checkvat").bootstrapSwitch('disabled', false)
    if (purchaseorder.names() != "Non Inventory") {
        $(".codeitem").on('click');
    }
}

purchaseorder.resetView = function () {
    ko.mapping.fromJS(purchaseorder.newRecord(), purchaseorder.record)
    $(".formInput").val("")
    $("#supliername").val("")
    $(".Amount").val("")
    $("#supliercode").val(0)
    $('.nav-tabs a[href="#Create"]').tab('show')
    purchaseorder.text("Create Purchase Order")
    $('#supliercode').data('kendoDropDownList').value(-1);
    $('#payment').data('kendoDropDownList').value(-1);
    $('#typepurchase').data('kendoDropDownList').value(-1);
    purchaseorder.enableView()
    // purchaseorder.record.DocumentNumber(purchaseorder.setDocumentNumber())
    purchaseorder.maskingMoney()
    purchaseorder.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    $("#checkvat").bootstrapSwitch('state', false)
    purchaseorder.TotalAll(0)
    purchaseorder.GrandTotalAll(0)
}

purchaseorder.createdForm = function () {
    purchaseorder.resetView()
    // purchaseorder.record.DocumentNumber(purchaseorder.setDocumentNumber())
    purchaseorder.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    purchaseorder.showCreate(true)
    purchaseorder.reset(true)
    purchaseorder.postDraft(false)
    purchaseorder.printPDF(false)
    purchaseorder.deleteDraft(false)
    purchaseorder.backToDraft(false)

    if (purchaseorder.names() == "Non Inventory") {
        $(".invhide").hide()
        $(".item").prop('disabled', false);

    } else {
        $(".invhide").show()
        $(".item").prop('disabled', true);
    }
}

purchaseorder.editDraft = function (e) {
    // console.log(e)
    purchaseorder.createdForm()
    var data = []
    if (purchaseorder.names() == "Non Inventory") {
        setTimeout(function () { $(".invhide").hide(); }, 50);
        $(".item").prop('disabled', false);
        data = _.find(purchaseorder.dataMasterPurchaseOrder(), function (o) {
            return o.ID == e;
        });
        //    console.log(data)
    } else {
        $(".invhide").show()
        $(".item").prop('disabled', true);
        data = _.find(purchaseorder.dataMasterPurchaseInventory(), function (o) {
            return o.ID == e;
        });
        //  console.log(data)
    }

    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    // data.Remark = data.Remark.substring(16)
    if (data.Currency == "IDR") {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceUSD = 0
            data.ListDetail[i].AmountUSD = 0
            purchaseorder.TotalAll(data.TotalIDR)
            purchaseorder.GrandTotalAll(data.GrandTotalIDR)
        }
    } else {
        for (var i = 0; i < data.ListDetail.length; i++) {
            data.ListDetail[i].PriceIDR = 0
            data.ListDetail[i].AmountIDR = 0
            purchaseorder.TotalAll(data.TotalUSD)
            purchaseorder.GrandTotalAll(data.GrandTotalUSD)
        }
    }


    // data.Currency = "TEST"
    ko.mapping.fromJS(data, purchaseorder.record)
    if (data.DownPayment == 0) {
        purchaseorder.record.DownPayment("0")
    }
    if (data.VAT > 0) {
        $('#checkvat').bootstrapSwitch('state', true);
    } else {
        $('#checkvat').bootstrapSwitch('state', false);
    }
    purchaseorder.record.ID(e)
    newArr = []
    _.each(purchaseorder.record.ListDetail(), function (v, i) {
        if (v != undefined) {
            newArr.push(purchaseorder.listDetail(ko.mapping.toJS(v)))
        }
    })

    purchaseorder.record.ListDetail(newArr)
    purchaseorder.reset(false)
    purchaseorder.showCreate(true)
    purchaseorder.postDraft(true)
    purchaseorder.backToDraft(true)
    purchaseorder.deleteDraft(true)
    purchaseorder.printPDF(false)
    $(".btnDelete").attr('disabled', false);
    $(".formInput").removeAttr("disabled")
    purchaseorder.text("Edit Draft Purchase Order")
    purchaseorder.enableView()
    purchaseorder.checkdata()
    purchaseorder.maskingMoney()
}

purchaseorder.backDraft = function () {
    $('.nav-tabs a[href="#List"]').tab('show')
    $("#List").addClass("active");
    $("#Create").removeClass("active");

}
purchaseorder.delete = function () {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + ko.mapping.toJS(purchaseorder.record).DocumentNumber + "?",
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
            var url = "/transaction/deletedraft";
            var param = {
                id: ko.mapping.toJS(purchaseorder.record).ID,
                documentNumber: ko.mapping.toJS(purchaseorder.record).DocumentNumber,
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
                            window.location.assign("/transaction/purchaseorder")
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

purchaseorder.onChangeStatus = function (textFilter) {
    purchaseorder.dataMasterPurchaseOrder([])
    var allData = purchaseorder.dataMasterPurchaseOrderOriginal()
    if (textFilter != "" || textFilter != undefined) {

        var Data = _.filter(allData, function (o) {
            return o.Status.indexOf(textFilter) > -1
        });
        purchaseorder.dataMasterPurchaseOrder(Data)
    }
    purchaseorder.renderGrid()

}
purchaseorder.toListTab = function(){
    purchaseorder.filterAlertType("tab")
    purchaseorder.search()
}
purchaseorder.search = function (e) {
    purchaseorder.filterAlert(true)
    if (purchaseorder.filterAlertType() == "tab"){
        purchaseorder.filterAlert(false)
    }
    if (purchaseorder.names() == "Non Inventory") {
        purchaseorder.textSearch(e)
        purchaseorder.filterindicator(true)
        purchaseorder.getDataPurchaseOrder(function () {
            purchaseorder.renderGrid()
        })
    } else {
        purchaseorder.textSearch(e)
        purchaseorder.filterindicatorInv(true)
        purchaseorder.getDataPurchaseInventory(function () {
            purchaseorder.renderGridInventory()
        })
    }
    purchaseorder.filterAlertType("search")
}

purchaseorder.choose = function (choosenote) {
    purchaseorder.filterAlert(false)
    $('.nav-tabs a[href="#List"]').tab('show')
    $('#textSearch').val("")
    if (choosenote == "ACTIVA") {
        purchaseorder.names("Non Inventory")
        setTimeout(function () {
            $(".invhide").hide()
            $(".typehide").show()
        }, 50);

        $(".item").prop('disabled', false);
        purchaseorder.getDataPurchaseOrder(function () {
            purchaseorder.renderGrid()
            $("#noninv").css({ "background-color": "#ffffff", "color": "#3ea49d", "border-bottom": "4px solid #f4222d" });
            $("#inv").css({ "background-color": "#45b6af", "color": "white", "border-color": "#3ea49d", "border-bottom": "1px solid #3ea49d" });
        })
    } else {
        purchaseorder.getDataPurchaseInventory(function () {
            purchaseorder.renderGridInventory()
        })
        purchaseorder.names("Inventory")
        setTimeout(function () {
            $(".invhide").show()
            $(".typehide").hide()
        }, 50);

        $(".item").prop('disabled', true);
        $("#inv").css({ "background-color": "#ffffff", "color": "#3ea49d", "border-bottom": "4px solid #f4222d" });
        $("#noninv").css({ "background-color": "#45b6af", "color": "white", "border-color": "#3ea49d", "border-bottom": "1px solid #3ea49d" });

    }
    // purchaseorder.resetView()
}

purchaseorder.getDataSales = function () {
    model.Processing(true)
    ajaxPost('/master/getdatasales', {}, function (res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseorder.dataMasterSales(res.Data)
        var DataSales = res.Data
        for (i in DataSales) {
            DataSales[i].Kode = DataSales[i].SalesID + ""
            DataSales[i].Name = DataSales[i].SalesID + " - " + DataSales[i].SalesName
        }
        purchaseorder.dataDropDownSales(DataSales)

        model.Processing(false)
    })
}

purchaseorder.getDataInventory = function () {
    model.Processing(true)

    ajaxPost('/master/getalldatainventorycentral', {}, function (res) {
        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseorder.dataMasterInventory(res.Data)
        var DataInventory = res.Data
        for (i in DataInventory) {
            DataInventory[i].Kode = DataInventory[i].INVID + ""
            DataInventory[i].Name = DataInventory[i].INVID + " - " + DataInventory[i].INVDesc
        }
        purchaseorder.dataDropDownInventory(DataInventory)

        model.Processing(false)
    })
}

purchaseorder.onChangeCodeItem = function (value, index) {
    findaccount = _.find(purchaseorder.dataMasterInventory(), {
        INVID: value
    })
    purchaseorder.record.ListDetail()[index].Item(findaccount.INVDesc);
}

purchaseorder.init = function () {
    ProActive.KendoDatePickerRange();
    ProActive.GlobalSearch("gridListPurchaseOrder", 
    ["SalesName", "Remark", "SupplierName", "DocumentNumber", "TotalIDR equals", "TotalUSD equals", "GrandTotalIDR equals", "GrandTotalUSD equals"]);
    purchaseorder.maskingMoney()
    purchaseorder.getDateNow()
    purchaseorder.switchButton()
    purchaseorder.setDate()
    purchaseorder.getDataSales()
    purchaseorder.getDataInventory()
    $("#noninv").css({ "background-color": "#ffffff", "color": "#3ea49d", "border-bottom": "4px solid #f4222d" });
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null || type != null) {
        if (type == "POINV") {
            purchaseorder.names("Inventory")
            purchaseorder.choose('IVENTORY')
        } else {
            purchaseorder.getDataPurchaseOrder(function () {
                purchaseorder.renderGrid()
            })
        }
    } else {
        purchaseorder.getDataPurchaseOrder(function () {
            purchaseorder.renderGrid()
        })
    }
}
purchaseorder.fromPOINVSummary = function () {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null) {
        var allData = purchaseorder.dataMasterPurchaseOrderOriginal()
        if (type != null) {
            if (type == "POINV") {
                purchaseorder.names("Inventory")
                // purchaseorder.choose('IVENTORY')
                allData = purchaseorder.dataMasterPurchaseInventoryOriginal()
                // console.log(purchaseorder.dataMasterPurchaseInventoryOriginal())
            }
        }

        var data = _.find(allData, function (o) {
            return o.DocumentNumber == num;
        });
        //  console.log(purchaseorder.dataMasterPurchaseInventory())
        if (data != undefined) {
            purchaseorder.viewDraft(data.ID)
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning",
                confirmButtonColor: "#3da09a"
            }, function () {
                window.location.assign("/transaction/purchaseorder")
            });
        }
    }
}

purchaseorder.detailReportPdf = function () {
    var valuefilter = ""
    if (purchaseorder.names() == "Non Inventory") {
        valuefilter = purchaseorder.textSupplierSearch()
    } else {
        valuefilter = purchaseorder.textSupplierSearch()
    }

    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: moment(startdate).format("YYYY-MM-DD"),
        DateEnd: moment(enddate).format("YYYY-MM-DD"),
        ReportType: "detail",
        ReportBy: "Customer",
        ValueFilter: valuefilter,
        IsInventory: purchaseorder.names() == "Inventory"
    }
    model.Processing(true)
    ajaxPost("/report/exporttopdfpurchaseorderdetail", param, function (res) {
        if (res.IsError) {
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res, '_blank');
        model.Processing(false)
    })
}

$(function () {
    purchaseorder.getDataSupplier()
    purchaseorder.getDataAccount()
    purchaseorder.getDataTypePurchase()
    purchaseorder.init()
    purchaseorder.text("Created Purchase Order")
    $(".priceusd").keyup(function () {
        var nameLength = $(".priceusd").val().length;
        if (nameLength > 0) {
            $(".priceidr").attr("disabled", "disabled");
        } else {
            $(".priceidr").removeAttr("disabled");
        }
    });
    $(".priceidr").keyup(function () {
        var nameLength = $(".priceidr").val().length;
        if (nameLength > 0) {
            $(".priceusd").attr("disabled", "disabled");
        } else {
            $(".priceusd").removeAttr("disabled");
        }
    });
    purchaseorder.names("Non Inventory")
});