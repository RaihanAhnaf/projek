var purchaseinvoice = {}
var processs = 0;
purchaseinvoice.date = ko.observable()
purchaseinvoice.dataMasterCustomer = ko.observableArray([])
purchaseinvoice.dataDropDownCustomer = ko.observableArray([])
purchaseinvoice.dataDropDownAccount = ko.observableArray([])
purchaseinvoice.dataDropDownPO = ko.observableArray([])
purchaseinvoice.dataDropDownSupplierFilter = ko.observableArray([])
purchaseinvoice.dataMasterPurchaseOrder = ko.observableArray([])
purchaseinvoice.dataMasterAccount = ko.observableArray([])
purchaseinvoice.dataMasterInvoice = ko.observableArray([])
purchaseinvoice.dataMasterInvoiceOriginal = ko.observableArray([])
purchaseinvoice.dataMasterInvoiceInventory = ko.observableArray([])
purchaseinvoice.dataMasterInvoiceOriginalInventory = ko.observableArray([])
// purchaseinvoice.dropDownPO = ko.observableArray([])
purchaseinvoice.GrandTotalAll = ko.observable(0)
purchaseinvoice.statusText = ko.observable()
purchaseinvoice.textSearch = ko.observable()
purchaseinvoice.text = ko.observable()
purchaseinvoice.DatePageBar = ko.observable()
purchaseinvoice.filterList = ko.observable(false)
purchaseinvoice.textSupplierSearch = ko.observable("")
purchaseinvoice.showCreate = ko.observable(false)
purchaseinvoice.backToList = ko.observable(false)
purchaseinvoice.showPrint = ko.observable(false)
purchaseinvoice.showAssetTable = ko.observable(true)
purchaseinvoice.valueDepartment = ko.observable("")
purchaseinvoice.dataMasterCategory = ko.observableArray([])
purchaseinvoice.monthDepreciation = ko.observable(false)
purchaseinvoice.names = ko.observable()
purchaseinvoice.dropDownPOMaster = ko.observableArray([])
purchaseinvoice.dropDownPOMasterInventory = ko.observableArray([])
purchaseinvoice.editingMode = ko.observable(false);
purchaseinvoice.filterAlert = ko.observable(false) 
purchaseinvoice.blurMD = function(value, index) {
    var classExist = $("#tdmonth" + index).hasClass("tdBackground");
    //console.log(value, index, classExist)
    if (classExist) {
        if (value != "0") {
            $("#tdmonth" + index).removeClass("tdBackground");
        }
    }
}
purchaseinvoice.acccode = [{
    value: 5100,
    text: "SALES"
}, {
    value: 5200,
    text: "REVENUE"
}]

purchaseinvoice.BoolVat = ko.observable()
purchaseinvoice.sequenceNumber = ko.observable()
purchaseinvoice.roman = [{
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
purchaseinvoice.newRecord = function() {
    var page = {
        ID: "",
        AccountCode: "",
        DocumentNumber: "",
        SupplierCode: "",
        SupplierName: "",
        DateStr: "",
        DatePosting: "",
        PoNumber: "",
        Type: "",
        ListDetail: [],
        Total: 0,
        TotalIDR: 0,
        TotalUSD: 0,
        VAT: 0,
        Discount: 0,
        GrandTotal: 0,
        GrandTotalIDR: 0,
        GrandTotalUSD: 0,
        Status: "",
        Remark: "",
        Payment: "",
        Currency: "",
        Rate: 0.0,
        Discount: 0.0,
        DownPayment: 0,
        SalesCode: "",
        SalesName: "",
        Department: "",
        INVCMI: false,
        LocationID: 0,
        LocationName: "",
        DateStrPI: "",
        DatePostingPI: "",
        DocumentNumberPI: ""
    }

    page.ListDetail.push(purchaseinvoice.listDetail({}))

    return page
}
purchaseinvoice.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    purchaseinvoice.DatePageBar(page)
}

purchaseinvoice.listDetail = function(data) {
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

     //Menghitung Total
     var data = ko.mapping.toJS(purchaseinvoice.record)
     if (purchaseinvoice.record.ListDetail()[0].PriceIDR > 0 || total > 0) {
         var alldataidr = ko.mapping.toJS(purchaseinvoice.record.ListDetail())
         var totalIDR = _.sumBy(alldataidr, function(v) {
             return v.AmountIDR
         })
         purchaseinvoice.record.TotalIDR(totalIDR)
         purchaseinvoice.record.Total(ChangeToRupiah(totalIDR))
         var tot = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(totalIDR)
         purchaseinvoice.record.GrandTotalIDR(tot)
         purchaseinvoice.record.GrandTotal(ChangeToRupiah(tot))
     } 
    })

    return x
}
purchaseinvoice.assetDepreciation = function() {
    var dataTmp = {
        Item: '',
        Qty: 0,
        PriceIDR: 0,
        AmountIDR: 0,
        DatePosting: "",
        Category: "",
        SumMonthDepreciation: 0
    }
    var x = ko.mapping.fromJS(dataTmp)

    // var x = ko.mapping.fromJS(dataTmp)
    return x
}
purchaseinvoice.dataAssetDepreciation = ko.mapping.fromJS([])
purchaseinvoice.record = ko.mapping.fromJS(purchaseinvoice.newRecord())


purchaseinvoice.onChangePoNumber = function(value) {
    //console.log(value)

    $("#my-checkbox").bootstrapSwitch('disabled', false)
    $("#cx-cmi").bootstrapSwitch('disabled', false)
    var data = _.filter(purchaseinvoice.dataDropDownPO(), {
        'DocumentNumber': value,
    })

    if (data.length == 0) {
        $('#cx-cmi').bootstrapSwitch('state', false);
        $("#cx-cmi").bootstrapSwitch('disabled', true);
        $(".item").prop('disabled', true);
        $('#departmenDropdown').data("kendoDropDownList").select(0);
        return
    }

    console.log(data)
    purchaseinvoice.record.Discount(data[0].Discount)
    purchaseinvoice.record.Rate(data[0].Rate)
    purchaseinvoice.record.Currency(data[0].Currency)
    purchaseinvoice.record.Payment(data[0].Payment)
    purchaseinvoice.record.DocumentNumber(data[0].DocumentNumber)
    purchaseinvoice.record.ID(data[0].ID)
    purchaseinvoice.record.DownPayment(data[0].DownPayment)
    purchaseinvoice.record.DateStr(data[0].DateStr)
    purchaseinvoice.record.Type(data[0].Type)
    purchaseinvoice.record.SupplierCode(data[0].SupplierCode)
    purchaseinvoice.record.SupplierName(data[0].SupplierName)
    purchaseinvoice.record.AccountCode(data[0].AccountCode)
    purchaseinvoice.record.Remark(data[0].Remark)
    purchaseinvoice.record.Total(ChangeToRupiah(data[0].TotalIDR))
    purchaseinvoice.record.TotalIDR(data[0].TotalIDR)
    purchaseinvoice.record.TotalUSD(data[0].TotalUSD)
    purchaseinvoice.record.Discount(data[0].Discount)
    purchaseinvoice.record.GrandTotal(ChangeToRupiah(data[0].GrandTotalIDR))
    purchaseinvoice.record.GrandTotalIDR(data[0].GrandTotalIDR)
    purchaseinvoice.record.GrandTotalUSD(data[0].GrandTotalUSD)
    purchaseinvoice.record.SalesCode(data[0].SalesCode)
    purchaseinvoice.record.SalesName(data[0].SalesName)
    purchaseinvoice.record.LocationID(data[0].LocationID)
    purchaseinvoice.record.LocationName(data[0].LocationName)
    purchaseinvoice.record.DatePosting(moment(data[0].DatePosting).format('DD-MMM-YYYY'))
    purchaseinvoice.record.DateStr(data[0].DateStr)
    purchaseinvoice.record.DatePostingPI(moment(new Date()).startOf('day').format('DD-MMM-YYYY'))
    purchaseinvoice.record.DocumentNumberPI("")
    $("#dateposting").data('kendoDatePicker').min(moment(data[0].DatePosting).startOf('day').toDate());

    if (data[0].VAT != 0) {
        $('#my-checkbox').bootstrapSwitch('state', true);
    } else {
        $('#my-checkbox').bootstrapSwitch('state', false);
    }
    $("#my-checkbox").bootstrapSwitch('disabled', true);

    $('#cx-cmi').bootstrapSwitch('state', data[0].INVCMI == true);
    $("#cx-cmi").bootstrapSwitch('disabled', false);

    if (data.length != 0) {
        purchaseinvoice.record.ListDetail([])
        purchaseinvoice.dataAssetDepreciation([])
        for (var i = 0; i < data[0].ListDetail.length; i++) {
            purchaseinvoice.record.ListDetail.push(purchaseinvoice.listDetail({}))
            purchaseinvoice.maskingMoney()

            purchaseinvoice.record.ListDetail()[i].Id(data[0].ListDetail[i].Id)
            purchaseinvoice.record.ListDetail()[i].CodeItem(data[0].ListDetail[i].CodeItem)
            purchaseinvoice.record.ListDetail()[i].Item(data[0].ListDetail[i].Item)
            purchaseinvoice.record.ListDetail()[i].Qty(data[0].ListDetail[i].Qty)
            purchaseinvoice.record.ListDetail()[i].PriceUSD(data[0].ListDetail[i].PriceUSD)
            purchaseinvoice.record.ListDetail()[i].PriceIDR(data[0].ListDetail[i].PriceIDR)
            purchaseinvoice.record.ListDetail()[i].AmountUSD(data[0].ListDetail[i].AmountUSD)
            purchaseinvoice.record.ListDetail()[i].AmountIDR(data[0].ListDetail[i].AmountIDR)
            if (data[0].Currency == "USD") {
                purchaseinvoice.record.ListDetail()[i].PriceIDR(0)
                purchaseinvoice.record.ListDetail()[i].AmountIDR(0)
            } else {
                purchaseinvoice.record.ListDetail()[i].PriceUSD(0)
                purchaseinvoice.record.ListDetail()[i].AmountUSD(0)
            }
            purchaseinvoice.dataAssetDepreciation.push(purchaseinvoice.assetDepreciation())
            //console.log(purchaseinvoice.dataAssetDepreciation()[i])
            purchaseinvoice.dataAssetDepreciation()[i].Item(data[0].ListDetail[i].Item)
            purchaseinvoice.dataAssetDepreciation()[i].Qty(data[0].ListDetail[i].Qty)
            purchaseinvoice.dataAssetDepreciation()[i].PriceIDR(data[0].ListDetail[i].PriceIDR)
            purchaseinvoice.dataAssetDepreciation()[i].AmountIDR(data[0].ListDetail[i].AmountIDR)
            purchaseinvoice.dataAssetDepreciation()[i].DatePosting(data[0].DatePosting)
        }
    } else {
        purchaseinvoice.record.ListDetail([])
        purchaseinvoice.record.ListDetail.push(purchaseinvoice.listDetail({}))
        purchaseinvoice.record.Total(0)
    }
    $(".Amount").attr("disabled", "disabled")
    if (purchaseinvoice.names() == "Non Inventory") {
        purchaseinvoice.names("Non Inventory")
        $(".noninv").show()
        $(".invhide").hide()
        $(".item").prop('disabled', false);
        $('#departmenDropdown').data("kendoDropDownList").select(0);
        purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMaster())
        $(".bootstrap-switch-id-cx-cmi").parent().parent().hide();

    } else {
        $('#departmenDropdown').data("kendoDropDownList").value("COMMERCE");
        purchaseinvoice.names("Inventory")
        $(".noninv").hide()
        $(".invhide").show()
        $(".item").prop('disabled', true);
        purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMasterInventory())
        $(".bootstrap-switch-id-cx-cmi").parent().parent().show();
    }
}

purchaseinvoice.switchButton = function() {
    $('#my-checkbox').on('switchChange.bootstrapSwitch', function(event, state) {
        var data = ko.mapping.toJS(purchaseinvoice.record)
        purchaseinvoice.BoolVat(state)
        // if (state) {
        //     var totVat = FormatCurrency(data.Total) / 10
        //     purchaseinvoice.record.VAT(ChangeToRupiah(totVat))
        // } else {
        //     purchaseinvoice.record.VAT(0)
        // }
        if (state) {
            if (purchaseinvoice.record.Currency() == "IDR") {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalIDR)
                purchaseinvoice.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
                var Totalall = GT + totVat
                purchaseinvoice.record.GrandTotalIDR(Totalall)
                purchaseinvoice.GrandTotalAll(ChangeToRupiah(Totalall))

            } else {
                var totVat = ((100 - FormatCurrency(data.Discount)) / 1000) * FormatCurrency(data.TotalUSD)
                purchaseinvoice.record.VAT(ChangeToRupiah(totVat))
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
                var Totalall = GT + totVat
                purchaseinvoice.record.GrandTotalUSD(Totalall)
                purchaseinvoice.record.GrandTotalIDR(Totalall * FormatCurrency(data.Rate))
                purchaseinvoice.GrandTotalAll(ChangeToRupiah(Totalall))

            }
        } else {
            purchaseinvoice.record.VAT(ChangeToRupiah(0))
            if (purchaseinvoice.record.Currency() == "IDR") {
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalIDR)
                purchaseinvoice.record.GrandTotalIDR(GT)
                purchaseinvoice.GrandTotalAll(ChangeToRupiah(GT))

            } else {
                var GT = ((100 - FormatCurrency(data.Discount)) / 100) * FormatCurrency(data.TotalUSD)
                purchaseinvoice.record.GrandTotalUSD(GT)
                purchaseinvoice.record.GrandTotalIDR(GT * FormatCurrency(data.Rate))
                purchaseinvoice.GrandTotalAll(ChangeToRupiah(GT))

            }
        }
    });
}



purchaseinvoice.checkdata = function() {
    var data = ko.mapping.toJS(purchaseinvoice.record)
    if (FormatCurrency(data.ListDetail[0].PriceIDR) > 0) {
        $(".priceusd").attr("disabled", "disabled");
        $(".priceidr").removeAttr("disabled");
    } else if (FormatCurrency(data.ListDetail[0].PriceUSD) > 0) {
        $(".priceidr").attr("disabled", "disabled");
        $(".priceusd").removeAttr("disabled");
    }
}

purchaseinvoice.GetDataPurchaseInvoice = function() {
    return new Promise(resolve =>{
        model.Processing(true)
        ajaxPost('/transaction/getdatapurchaseinvoice', {}, function(res) {
            if (res.IsError === "true") {
                swal({
                    title: "Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor: "#3da09a"
                })
                return
            }
            purchaseinvoice.dropDownPOMaster(res.Data)
            // purchaseinvoice.dataDropDownPO(res.Data)
            purchaseinvoice.processsPlusPlus()
            model.Processing(false)
            resolve(res.Data)
        })
    })
}

purchaseinvoice.GetDataPurchaseInvoiceInventory = function() {
    return new Promise(resolve =>{
        model.Processing(true)
        ajaxPost('/transaction/getdatapurchaseinvoiceinventory', {}, function(res) {
            if (res.IsError === "true") {
                swal({
                    title: "Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor: "#3da09a"
                })
                return
            }
            purchaseinvoice.dropDownPOMasterInventory(res.Data)
            // purchaseinvoice.dataDropDownPO(res.Data)
            purchaseinvoice.processsPlusPlus()
            model.Processing(false)
            resolve(res.Data)
        })
    })
    
}

purchaseinvoice.getDataPurchaseInvoiceStatusPI = function(callback) {
    var param = {}
    if (purchaseinvoice.filterList() == true) {
        param = {
            DateStart: moment($("#dateStart").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            DateEnd: moment($("#dateEnd").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            Filter: purchaseinvoice.filterList(),
            TextSearch: purchaseinvoice.textSearch(),
            SupplierCode: purchaseinvoice.textSupplierSearch(),
        }
    } else {
        param = {
            DateStart: moment($("#dateStart").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            DateEnd: moment($("#dateEnd").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            Filter: purchaseinvoice.filterList(),
        }
    }
    //console.log(param)
    ajaxPost('/transaction/getdatapurchaseinvoicestatuspi', param, function(res) {
        if (res.IsError && purchaseinvoice.filterAlert()) {
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
        //console.log(res.Data)
        purchaseinvoice.dataMasterInvoice(res.Data)
        purchaseinvoice.dataMasterInvoiceOriginal(res.Data)
        purchaseinvoice.processsPlusPlus()
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
purchaseinvoice.getDataPurchaseInvoiceStatusPIInventory = function(callback) {
    var param = {}
    if (purchaseinvoice.filterList() == true) {
        param = {
            DateStart: moment($("#dateStart").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            DateEnd: moment($("#dateEnd").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            Filter: purchaseinvoice.filterList(),
            TextSearch: purchaseinvoice.textSearch(),
            SupplierCode: purchaseinvoice.textSupplierSearch(),
        }
    } else {
        param = {
            DateStart: moment($("#dateStart").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            DateEnd: moment($("#dateEnd").data("kendoDatePicker").value()).startOf('day').format("YYYY-MM-DD") + "T00:00:00Z",
            Filter: purchaseinvoice.filterList(),
        }
    }
    ajaxPost('/transaction/getdatapurchaseinvoicestatuspiinventory', param, function(res) {
        if (res.IsError&& purchaseinvoice.filterAlert()) {
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
        purchaseinvoice.dataMasterInvoiceInventory(res.Data)
        purchaseinvoice.dataMasterInvoiceOriginalInventory(res.Data)
        purchaseinvoice.processsPlusPlus()
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

purchaseinvoice.getDataAccount = function() {
    model.Processing(true)

    ajaxPost('/transaction/getaccount', {}, function(res) {

        if (res.Total === 0) {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseinvoice.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        purchaseinvoice.dataDropDownAccount(DataAccount)
        purchaseinvoice.processsPlusPlus()
        model.Processing(false)
    })
}

purchaseinvoice.addNewItem = function() {
    purchaseinvoice.record.ListDetail.push(purchaseinvoice.listDetail({}))
    purchaseinvoice.maskingMoney()
    purchaseinvoice.checkdata()
    if (purchaseinvoice.names() == "Non Inventory") {
        $(".invhide").hide()
        $(".item").prop('disabled', false);

    } else {
        $(".invhide").show()
        $(".item").prop('disabled', true);
    }
}
purchaseinvoice.formCreated = function() {

    //asset decition
    purchaseinvoice.showAssetTable(true)
    if ($('#tableListDetail').attr('class') == "col-md-12") {
        $("#tableListDetail").toggleClass('col-md-12 col-md-9');
    }
    //end
    $("#my-checkbox").bootstrapSwitch('disabled', false)
    $("#cx-cmi").bootstrapSwitch('disabled', false)
    purchaseinvoice.resetView()
    purchaseinvoice.showCreate(true)
    purchaseinvoice.showPrint(false)
    purchaseinvoice.backToList(false)
    purchaseinvoice.enableView()
    purchaseinvoice.record.ListDetail([])
    purchaseinvoice.dataAssetDepreciation([])
    purchaseinvoice.record.ListDetail.push(purchaseinvoice.listDetail({}))
    purchaseinvoice.dataAssetDepreciation.push(purchaseinvoice.assetDepreciation())
    purchaseinvoice.maskingMoney()
    purchaseinvoice.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    purchaseinvoice.dataDropDownPO([])
    // purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMaster())
    $("#my-checkbox").bootstrapSwitch('state', false)
    $("#cx-cmi").bootstrapSwitch('state', false)
    $("#labelPurchaseNo").hide()
    $("#purchasenodropdown").removeAttr("style")
    $(".Amount").attr("disabled", "disabled")
    $("#my-checkbox").bootstrapSwitch('disabled', true)
    $("#cx-cmi").bootstrapSwitch('disabled', true)
    $('#departmenDropdown').data("kendoDropDownList").enable(true);

    if (purchaseinvoice.names() == "Non Inventory") {
        purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMaster())
        $(".noninv").show()
        $(".invhide").hide()
        $(".item").prop('disabled', false);
        $('#departmenDropdown').data("kendoDropDownList").select(0); 
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
         $('#tableListDetail').removeClass('col-md-12').addClass('col-md-9');
         
        $(".bootstrap-switch-id-cx-cmi").parent().parent().hide();

    } else {
        purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMasterInventory())
        $('#departmenDropdown').data("kendoDropDownList").value("COMMERCE");
        purchaseinvoice.names("Inventory")
        $(".noninv").hide()
        $(".invhide").show()
        $(".item").prop('disabled', true);
        setTimeout(function(){ $("#inv").css({
            "background-color": "#ffffff",
            "color": "#3ea49d",
            "border-bottom": "4px solid #f4222d"
        });
        $("#noninv").css({
            "background-color": "#45b6af",
            "color": "white",
            "border-color": "#3ea49d",
            "border-bottom": "1px solid #3ea49d"
        }); }, 50);
         $('#tableListDetail').removeClass('col-md-9').addClass('col-md-12');
         
        $(".bootstrap-switch-id-cx-cmi").parent().parent().show();
    }
    $("#dateposting").data("kendoDatePicker").enable(true);
    purchaseinvoice.searchOrder();
    purchaseinvoice.editingMode(false);
}

purchaseinvoice.monthRoman = function(e) {
    var data = _.find(purchaseinvoice.roman, e);
    return _.values(data)
}

purchaseinvoice.reset = function() {
    purchaseinvoice.record.ListDetail([]);
    purchaseinvoice.record.Total(0)
    purchaseinvoice.resetView()
    purchaseinvoice.editingMode(false);
}

purchaseinvoice.removerow = function() {
    purchaseinvoice.checkdata()
    purchaseinvoice.record.ListDetail.remove(this)
    if (purchaseinvoice.record.ListDetail().length == 0) {
        purchaseinvoice.record.ListDetail.push(purchaseinvoice.listDetail({}))
        purchaseinvoice.record.Total(0)
    } else {
        //Menghitung Total 
        var alldataidr = ko.mapping.toJS(purchaseinvoice.record.ListDetail())
        var totalIDR = _.sumBy(alldataidr, function(v) {
            return v.AmountIDR
        })

        var alldata = ko.mapping.toJS(purchaseinvoice.record.ListDetail())
        var totalUSD = _.sumBy(alldata, function(v) {
            return v.AmountUSD
        })

        if (FormatCurrency(purchaseinvoice.record.ListDetail()[0].PriceIDR()) > 0) {
            purchaseinvoice.record.Total(ChangeToRupiah(totalIDR))
        }

        if (FormatCurrency(purchaseinvoice.record.ListDetail()[0].PriceUSD()) > 0) {
            purchaseinvoice.record.Total(ChangeToRupiah(totalUSD))
        }
    }
}

purchaseinvoice.maskingMoney = function() {
    $('.currency').inputmask("numeric", {
        radiresultPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        // prefiresult: '$', //No Space, this will truncate the first character
        rightAlign: false,
        // oncleared: function () { self.Value(''); }
    });
}


purchaseinvoice.renderGrid = function() {
    var mydata = purchaseinvoice.dataMasterInvoice()
    var columns = [{
        title: 'Action',
        width: 55,
        template: "<button onclick='purchaseinvoice.viewDraft(\"#: ID #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>"
    }, {
        field: 'DateStrPI',
        title: 'Date Posting',
        width: 160,
    }, {
        field: 'DocumentNumberPI',
        title: 'Doc Number',
        width: 160,
    }, {
        field: 'SupplierName',
        title: 'Supplier Name',
        width: 200,
    }, {
        field: 'GrandTotalIDR',
        title: 'Order Total',
        template: "#=ChangeToRupiah(GrandTotalIDR)#",
        attributes: {
            "class": "rightAction",
        },
        // template: "#if ( Currency == 'USD') {# $ #: ChangeToRupiah(GrandTotalUSD) # #} else {# Rp. #: ChangeToRupiah(GrandTotalIDR) # #}#",
        width: 170,
    }, {
        field: 'Remark',
        title: 'Remark',
        width: 200,
    }]

    $('#gridListInvoice').kendoGrid({
        dataSource: {
            data: mydata,
            sort: {
                field: 'DatePosting',
                dir: 'desc',
            }
        },
        height: 500,
        width: 140,
        filterable: true,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "PurchaseInvoice", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
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

purchaseinvoice.renderGridInventory = function() {
    var mydata = purchaseinvoice.dataMasterInvoiceInventory()
    var columns = [{
        title: 'Action',
        width: 55,
        template: "<button onclick='purchaseinvoice.viewDraft(\"#: ID #\", \"inv\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>"
    }, {
        field: 'DateStrPI',
        title: 'Date Posting',
        width: 160,
    }, {
        field: 'DocumentNumberPI',
        title: 'Doc Number',
        width: 160,
    }, {
        field: 'SupplierName',
        title: 'Supplier Name',
        width: 200,
    }, {
        field: 'GrandTotalIDR',
        title: 'Order Total',
        template: "#=ChangeToRupiah(GrandTotalIDR)#",
        attributes: {
            "class": "rightAction",
        },
        // template: "#if ( Currency == 'USD') {# $ #: ChangeToRupiah(GrandTotalUSD) # #} else {# Rp. #: ChangeToRupiah(GrandTotalIDR) # #}#",
        width: 170,
    }, {
        field: 'Remark',
        title: 'Remark',
        width: 200,
    }]

    $('#gridListInvoice').kendoGrid({
        dataSource: {
            data: mydata,
            sort: {
                field: 'DatePosting',
                dir: 'desc',
            }
        },
        height: 500,
        width: 140,
        filterable: true,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "PurchaseInvoiceInventory", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
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
purchaseinvoice.exportExcel = function () {
    $("#gridListInvoice").getKendoGrid().saveAsExcel();
}
purchaseinvoice.changeFuncCategory = function(value, index) {
    if (value != "") {
        $("#tdcategory" + index).removeClass("tdBackground");
        $("#tdcategory" + index + ">span").removeClass("tdBackground");
    }
}
purchaseinvoice.saveData = function() {
    var data = ko.mapping.toJS(purchaseinvoice.record)

    data.DatePosting = moment($('#dateinvoice').data('kendoDatePicker').value()).startOf("day").format("YYYY-MM-DD") + "T00:00:00Z";
    data.DatePostingPI = moment($('#dateposting').data('kendoDatePicker').value()).startOf("day").format("YYYY-MM-DD") + "T00:00:00Z";
    data.DateStrPI = moment($('#dateposting').data('kendoDatePicker').value()).startOf("day").format("DD-MMM-YYYY");
    data.Status = "PI"
    data.AccountCode = parseInt(data.AccountCode)
    data.Department = purchaseinvoice.valueDepartment()
    data.DownPayment = FormatCurrency(data.DownPayment)
    // data.Total = FormatCurrency(data.Total)
    data.VAT = FormatCurrency(data.VAT)
    data.INVCMI = $('#cx-cmi').bootstrapSwitch('state');
    if (data.INVCMI == true)
        data.Status = "PAID"
    // data.GrandTotal = FormatCurrency(data.GrandTotal)
    for (var i = 0; i < data.ListDetail.length; i++) {
        data.ListDetail[i].Qty = parseInt(data.ListDetail[i].Qty)
        data.ListDetail[i].PriceUSD = FormatCurrency(data.ListDetail[i].PriceUSD)
        data.ListDetail[i].PriceIDR = FormatCurrency(data.ListDetail[i].PriceIDR)
        data.ListDetail[i].AmountIDR = FormatCurrency(data.ListDetail[i].AmountIDR)
        data.ListDetail[i].AmountUSD = FormatCurrency(data.ListDetail[i].AmountUSD)
    }
    //asset depreciation data 
    var asset = ko.mapping.toJS(purchaseinvoice.dataAssetDepreciation())
    for (i in asset) {
        $("#tdmonth" + i).removeClass("tdBackground");
        $("#tdcategory" + i).removeClass("tdBackground");
        $("#tdcategory" + i + ">span").removeClass("tdBackground");
        asset[i].SumMonthDepreciation = parseInt(asset[i].SumMonthDepreciation)
        if (asset[i].Category == "" && asset[i].SumMonthDepreciation != 0) {
            $("#tdcategory" + i).addClass("tdBackground");
            $("#tdcategory" + i + ">span").addClass("tdBackground");
            return swal({
                title: "Warning!",
                text: "None Category Asset has selected",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        }
        if (asset[i].Category != "" && asset[i].SumMonthDepreciation == 0) {
            $("#tdmonth" + i).addClass("tdBackground");
            return swal({
                title: "Warning!",
                text: "Please check your summary month depreciation",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        }
    }
    //end
    if ($('#departmenDropdown').data("kendoDropDownList").value() == "") {
        return swal({
            title: "Warning!",
            text: "None Department has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.Total == 0) {
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
    } else if (data.AccountCode == "") {
        return swal({
            title: "Warning!",
            text: "No Account Number has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else {
        for (var i = 0; i < data.ListDetail.length; i++) {
            if (!(data.ListDetail[i].AmountIDR || data.ListDetail[i].AmountUSD)) {
                return swal({
                    title: "Warning!",
                    text: "Data in ListDetail line hasn't completed yet",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            } else {
                var param = {
                    Data: data,
                    Asset: asset,
                    Department: purchaseinvoice.valueDepartment(),
                }

                var url = "/transaction/savepurchaseinvoice"
                swal({
                    title: "Are you sure?",
                    text: "You will save this Purchase Invoice Non Inventory",
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
                                    purchaseinvoice.resetView()
                                    purchaseinvoice.GetDataPurchaseInvoice().then(resolve =>{
                                        $("#purchasenumber").data('kendoDropDownList').setDataSource(purchaseinvoice.dropDownPOMaster())
                                    })
                                    purchaseinvoice.getDataPurchaseInvoiceStatusPI(function() {
                                        purchaseinvoice.renderGrid()
                                    })
                                    
                                    purchaseinvoice.formCreated();
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
purchaseinvoice.saveDataInventory = function() {
    var data = ko.mapping.toJS(purchaseinvoice.record)
    data.DatePosting = moment($('#dateinvoice').data('kendoDatePicker').value()).startOf("day").format("YYYY-MM-DD") + "T00:00:00Z";
    data.DatePostingPI = moment($('#dateposting').data('kendoDatePicker').value()).startOf("day").format("YYYY-MM-DD") + "T00:00:00Z";
    data.DateStrPI = moment($('#dateposting').data('kendoDatePicker').value()).startOf("day").format("DD-MMM-YYYY");
    data.Status = "PI"
    // data.AccountCode = 1401
    data.Department = purchaseinvoice.valueDepartment()
    data.DownPayment = FormatCurrency(data.DownPayment)
    // data.Total = FormatCurrency(data.Total)
    data.VAT = FormatCurrency(data.VAT)
    data.INVCMI = $('#cx-cmi').bootstrapSwitch('state');
    if (data.INVCMI == true)
        data.Status = "PAID"
    // data.GrandTotal = FormatCurrency(data.GrandTotal)
    for (var i = 0; i < data.ListDetail.length; i++) {
        data.ListDetail[i].Qty = parseInt(data.ListDetail[i].Qty)
        data.ListDetail[i].PriceUSD = FormatCurrency(data.ListDetail[i].PriceUSD)
        data.ListDetail[i].PriceIDR = FormatCurrency(data.ListDetail[i].PriceIDR)
        data.ListDetail[i].AmountIDR = FormatCurrency(data.ListDetail[i].AmountIDR)
        data.ListDetail[i].AmountUSD = FormatCurrency(data.ListDetail[i].AmountUSD)
    }
    //asset depreciation data 
    var asset = ko.mapping.toJS(purchaseinvoice.dataAssetDepreciation())
    for (i in asset) {
        $("#tdmonth" + i).removeClass("tdBackground");
        $("#tdcategory" + i).removeClass("tdBackground");
        $("#tdcategory" + i + ">span").removeClass("tdBackground");
        asset[i].SumMonthDepreciation = parseInt(asset[i].SumMonthDepreciation)
        if (asset[i].Category == "" && asset[i].SumMonthDepreciation != 0) {
            $("#tdcategory" + i).addClass("tdBackground");
            $("#tdcategory" + i + ">span").addClass("tdBackground");
            return swal({
                title: "Warning!",
                text: "None Category Asset has selected",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        }
        if (asset[i].Category != "" && asset[i].SumMonthDepreciation == 0) {
            $("#tdmonth" + i).addClass("tdBackground");
            return swal({
                title: "Warning!",
                text: "Please check your summary month depreciation",
                type: "info",
                confirmButtonColor: "#3da09a"
            })
        }
    }
    //end
    if ($('#departmenDropdown').data("kendoDropDownList").value() == "") {
        return swal({
            title: "Warning!",
            text: "None Department has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else if (data.Total == 0) {
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
        for (var i = 0; i < data.ListDetail.length; i++) {
            if (!(data.ListDetail[i].AmountIDR || data.ListDetail[i].AmountUSD)) {
                return swal({
                    title: "Warning!",
                    text: "Data in ListDetail line hasn't completed yet",
                    type: "info",
                    confirmButtonColor: "#3da09a"
                })
            } else {
                var param = {
                    Data: data,
                    Asset: asset,
                    Department: purchaseinvoice.valueDepartment(),
                }

                var url = "/transaction/savepurchaseinvoiceinventory"
                swal({
                    title: "Are you sure?",
                    text: "You will save this Purchase Invoice Inventory",
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
                                     purchaseinvoice.resetView()
                                    // window.location.assign("/transaction/purchaseinvoice")
                                    // $("#btnSave").attr("disabled", "disabled")
                                    purchaseinvoice.GetDataPurchaseInvoiceInventory().then(resolve =>{
                                        $("#purchasenumber").data('kendoDropDownList').setDataSource(purchaseinvoice.dropDownPOMasterInventory())
                                    })
                                    purchaseinvoice.getDataPurchaseInvoiceStatusPIInventory(function() {
                                        purchaseinvoice.renderGridInventory()
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
        }
    }
}

purchaseinvoice.getDataSupplier = function() {
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
        var DataSupplier = res.Data
        for (i in DataSupplier) {
            DataSupplier[i].Kode = DataSupplier[i].Kode + ""
            DataSupplier[i].Name = DataSupplier[i].Kode + " - " + DataSupplier[i].Name
        }
        purchaseinvoice.dataDropDownSupplierFilter(DataSupplier)
        purchaseinvoice.processsPlusPlus()
        model.Processing(false)
    })
}
purchaseinvoice.onChangeStatus = function(textFilter) {
    purchaseinvoice.dataMasterInvoice([])
    var allData = purchaseinvoice.dataMasterInvoiceOriginal()
    if (textFilter != "" || textFilter != undefined) {

        var Data = _.filter(allData, function(o) {
            return o.Status.indexOf(textFilter) > -1
        });
        purchaseinvoice.dataMasterInvoice(Data)
    }
    purchaseinvoice.renderGrid()
}

purchaseinvoice.setDate = function() {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}

purchaseinvoice.search = function(e) {
    console.log(e)
    purchaseinvoice.filterAlert(true)
    if (e == "tab"){
        purchaseinvoice.filterAlert(false)
    }
    // purchaseinvoice.textSearch(e)
    purchaseinvoice.filterList(true)
    if (purchaseinvoice.names() == "Inventory") {
        purchaseinvoice.getDataPurchaseInvoiceStatusPIInventory(function() {
            purchaseinvoice.renderGridInventory()
        })
    }
    else {
        purchaseinvoice.getDataPurchaseInvoiceStatusPI(function() {
            purchaseinvoice.renderGrid()
        })
    }
}

purchaseinvoice.viewDraft = function(e, type) {
    
    //asset decition
    purchaseinvoice.showAssetTable(false)
    if ($('#tableListDetail').attr('class') == "col-md-9") {
        $("#tableListDetail").toggleClass('col-md-9 col-md-12');
    }
    //end
    $("#my-checkbox").bootstrapSwitch('disabled', false)
    $("#cx-cmi").bootstrapSwitch('disabled', false)
    var allData = purchaseinvoice.dataMasterInvoiceOriginal()
    if (type== "inv"){
        allData = purchaseinvoice.dataMasterInvoiceOriginalInventory()
    }
    var data = _.find(allData, function(o) {
        return o.ID == e;
    });
    // console.log(e, data)
    if (data.VAT > 0) {
        $('#my-checkbox').bootstrapSwitch('state', true);
    } else {
        $('#my-checkbox').bootstrapSwitch('state', false);
    }
    $("#my-checkbox").bootstrapSwitch('disabled', true)
    $('#cx-cmi').bootstrapSwitch('state', data.INVCMI == true);
    $("#cx-cmi").bootstrapSwitch('disabled', true);
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    data.DatePostingPI = moment(data.DatePostingPI).format('DD-MMM-YYYY')
    ko.mapping.fromJS(data, purchaseinvoice.record)
    // $("#purchasenumber").data("kendoDropDownList").text(data.DocumentNumber)
    // $("#purchasenumber").data("kendoDropDownList").value(data.DocumentNumber)
    purchaseinvoice.record.Total(ChangeToRupiah(data.TotalIDR))
    purchaseinvoice.record.GrandTotal(ChangeToRupiah(data.GrandTotalIDR))
    $('#departmenDropdown').data("kendoDropDownList").value(data.Department);
    // console.log(data.DocumentNumber)
    $("#labelPurchaseNo").removeAttr("style")
    $("#purchasenodropdown").hide()
    $("#listinvoice").removeClass("active")
    $("#createinvoice").addClass("active")
    $("#dateposting").data("kendoDatePicker").enable(false);
    purchaseinvoice.text("View Purchase Invoice")
    purchaseinvoice.showCreate(false)
    purchaseinvoice.showPrint(true)
    purchaseinvoice.backToList(true)
    purchaseinvoice.disableView()
    purchaseinvoice.maskingMoney()
    $("#cx-cmi").bootstrapSwitch('disabled', true);
    purchaseinvoice.editingMode(true);
    $('.nav-tabs a[href="#createinvoice"]').tab('show');
}

purchaseinvoice.backList = function() {
    $('.nav-tabs a[href="#listinvoice"]').tab('show')
    $("#listinvoice").addClass("active");
    $("#createinvoice").removeClass("active");
}

purchaseinvoice.disableView = function() {
    $(".formInput").attr("disabled", "disabled")
    $(".btnDeleteSummary").attr("disabled", "disabled")
    $(".Amount").attr("disabled", "disabled")
    $("#addnewitem").attr("disabled", "disabled")
    $("#downPayment").attr("disabled", "disabled")
    // $("#my-checkbox").bootstrapSwitch('toggleDisabsled',true,true)
    // $("#my-checkbox").bootstrapSwitch('disabled', true)
    $('#dateinvoice').data('kendoDatePicker').enable(false);
    var dropDown1 = $("#purchasenumber").data("kendoDropDownList");
    dropDown1.enable(false);
    var dropDown2 = $("#departmenDropdown").data("kendoDropDownList");
    dropDown2.enable(false);
    var dropDown3 = $("#accountcode").data("kendoDropDownList");
    dropDown3.enable(false);
}

purchaseinvoice.enableView = function() {
    $(".formInput").removeAttr("disabled")
    $(".btnDeleteSummary").removeAttr("disabled")
    $(".Amount").removeAttr("disabled")
    $("#addnewitem").removeAttr("disabled")
    // $("#my-checkbox").bootstrapSwitch('disabled', false)
    // $("#my-checkbox").bootstrapSwitch('toggleDisabled',true,true)
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#dateinvoice').data('kendoDatePicker').enable(condition);
    var dropDown1 = $("#purchasenumber").data("kendoDropDownList");
    dropDown1.enable(true);
    var dropDown3 = $("#accountcode").data("kendoDropDownList");
    dropDown3.enable(true);
}

purchaseinvoice.resetView = function() {
    ko.mapping.fromJS(purchaseinvoice.newRecord(), purchaseinvoice.record)
    ko.mapping.fromJS(purchaseinvoice.assetDepreciation(), purchaseinvoice.dataAssetDepreciation)
    $(".formInput").val("")
    purchaseinvoice.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    $("#customername").val("")
    $(".Amount").val("")
    purchaseinvoice.text("Create Purchase Invoice")
    $('#purchasenumber').data('kendoDropDownList').value(-1);
    $('#accountcode').data('kendoDropDownList').value(-1);
    purchaseinvoice.enableView()
}

purchaseinvoice.printToPdf = function() {
    model.Processing(true)
    var param = {
        Id: purchaseinvoice.record.ID(),
    }
    var url = "/transaction/exporttopdfpurchaseinvoice"
    if (purchaseinvoice.names() == "Inventory"){
        url = "/transaction/exporttopdfpurchaseinvoiceinventory"
    }
    ajaxPost(url, param, function(e) {
        model.Processing(false)
        var taborWindow = window.open('/res/docs/purchaseinvoice/' + e, '_blank');
        taborWindow.focus();
    })
}
purchaseinvoice.processsPlusPlus = function() {
    processs += 1
    //console.log(processs)
    if (processs >= 6) {
        setTimeout(function() {
            purchaseinvoice.fromPOINVSummary()
        }, 300)
        // setTimeut(purchasepayment.createPPFromPOInvSummary(),10000) 
    }
}

purchaseinvoice.departmentDropdown = function() {
    $("#departmenDropdown").html("")
    var data = []
    ajaxPost("/transaction/getdatadepartment", {}, function(res) {
        $("#departmenDropdown").kendoDropDownList({
            filter: "contains",
            dataTextField: "DepartmentName",
            dataValueField: "DepartmentName",
            dataSource: res.Data,
            optionLabel: 'Select one',
            change: function(e) {
                var dataitem = this.dataItem();
                purchaseinvoice.valueDepartment(dataitem.DepartmentName)
            }
        });
    })
}
purchaseinvoice.getDataCategory = function() {
    var param = {}
    ajaxPost('/master/getdatacategory', param, function(res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchaseinvoice.dataMasterCategory(res.Data)
    })
}

purchaseinvoice.choose = function(choosenote) {
    purchaseinvoice.filterAlert(false)
    // purchaseinvoice.formCreated()
    purchaseinvoice.dataDropDownPO([])
    $('#textSearch').val("")
    purchaseinvoice.textSearch("")
    if (choosenote == "ACTIVA") {
        purchaseinvoice.names("Non Inventory")
        $(".noninv").show()
        $(".invhide").hide()
        $(".item").prop('disabled', false);
        $('#departmenDropdown').data("kendoDropDownList").select(0);
        purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMaster())
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
        $('#tableListDetail').removeClass('col-md-12').addClass('col-md-9');
           purchaseinvoice.getDataPurchaseInvoiceStatusPI(function() {
            purchaseinvoice.renderGrid()
              //console.log("By click masuk kesini")
        })

        $(".bootstrap-switch-id-cx-cmi").parent().parent().hide();
  
    } else {
        purchaseinvoice.getDataPurchaseInvoiceStatusPIInventory(function() {
            purchaseinvoice.renderGridInventory()
        })
        //console.log($('#departmenDropdown').data("kendoDropDownList"))
        $('#departmenDropdown').data("kendoDropDownList").value("COMMERCE");
        purchaseinvoice.names("Inventory")
        $(".noninv").hide()
        $(".invhide").show()
        $(".item").prop('disabled', true);
        purchaseinvoice.dataDropDownPO(purchaseinvoice.dropDownPOMasterInventory())
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
        $('#tableListDetail').removeClass('col-md-9').addClass('col-md-12');
        $(".bootstrap-switch-id-cx-cmi").parent().parent().show();
    }
    purchaseinvoice.reset()
    $(".formInput").attr("disabled", "disabled")
    $('.nav-tabs a[href="#listinvoice"]').tab('show')
}

 purchaseinvoice.saveSwitch = function() {
    if (purchaseinvoice.names() == "Non Inventory") {        
         purchaseinvoice.saveData()
     } else {
         purchaseinvoice.saveDataInventory()      
     }
 }
 purchaseinvoice.fromPOINVSummary = function() {
    //console.log("tet")
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null) {
        var allData = purchaseinvoice.dataMasterInvoiceOriginal()
        if (type != null){
            if (type =="POINV"){
                purchaseinvoice.names("Inventory")
                allData = purchaseinvoice.dataMasterInvoiceOriginalInventory()
            }
        }
        var data = _.find(allData, function(o) {
            return o.DocumentNumber == num;
        });
        // console.log(num, data)
        if (data != undefined) {
            var typedraft = "po" 
            if (type =="POINV"){
                typedraft = "inv"
            }
            purchaseinvoice.viewDraft(data.ID, typedraft)
        } else {
            swal({
                title: "Warning!",
                text: "Data is not found",
                type: "warning",
                confirmButtonColor: "#3da09a"
            }, function() {
                window.location.assign("/transaction/purchaseinvoice")
            });
        }
    }
}

purchaseinvoice.searchOrder = function() {
    var name = purchaseinvoice.names();
    var startdate = $('#dateStartE').data('kendoDatePicker').value();
    var enddate = $('#dateEndE').data('kendoDatePicker').value();

    var filter = {
        DateStart: moment(startdate).format('YYYY-MM-DD'),
        DateEnd: moment(enddate).format('YYYY-MM-DD'),
        Filter: true,
    }

    var url = '/transaction/getdatapurchaseorder';
    if (name == "Inventory") {
        url = '/transaction/getdatapurchaseinventory';
    };

    model.Processing(true)
    ajaxPost(url, filter, function (res) {
        model.Processing(false)
        
        var columns = [{
            title: 'Action',
            width: 55,
            template: "<button onclick='purchaseinvoice.createPI(\"#: DocumentNumber #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil' aria-hidden='true'></i></button>"
        }, {
            field: 'DateStr',
            title: 'Date Posting',
            width: 160,
        }, {
            field: 'DocumentNumber',
            title: 'Doc Number',
            width: 160,
        }, {
            field: 'SupplierName',
            title: 'Supplier Name',
            width: 200,
        }, {
            field: 'GrandTotalIDR',
            title: 'Order Total',
            template: "#=ChangeToRupiah(GrandTotalIDR)#",
            attributes: {
                "class": "rightAction",
            },
            // template: "#if ( Currency == 'USD') {# $ #: ChangeToRupiah(GrandTotalUSD) # #} else {# Rp. #: ChangeToRupiah(GrandTotalIDR) # #}#",
            width: 170,
        }, {
            field: 'Remark',
            title: 'Remark',
            width: 200,
        }]
    
        $('#gridEditInvoice').html("");
        $('#gridEditInvoice').kendoGrid({
            dataSource: {
                data: res.Data,
                sort: {
                    field: 'DatePosting',
                    dir: 'desc',
                }
            },
            height: 500,
            width: 140,
            filterable: true,
            scrollable: true,
            columns: columns,
        })

    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
        model.Processing(false)
    })
}

purchaseinvoice.createPI = function(id) {
    purchaseinvoice.formCreated();
    purchaseinvoice.editingMode(true);
    //purchaseinvoice.record.DocumentNumber(id);
    purchaseinvoice.onChangePoNumber(id);
}

purchaseinvoice.init = function() {        
    ProActive.KendoDatePickerRange();     
    ProActive.KendoDatePickerRange("dateStartE", "dateEndE");
    ProActive.GlobalSearch("gridListInvoice", 
    ["SalesName", "Remark", "SupplierName", "DocumentNumber", "TotalIDR equals", "TotalUSD equals", "GrandTotalIDR equals", "GrandTotalUSD equals"]);

    var now = new Date()
    purchaseinvoice.date(moment(now).format("DD MMMM YYYY"))
    purchaseinvoice.setDate()
    purchaseinvoice.getDataAccount()
    purchaseinvoice.getDataCategory()
    purchaseinvoice.text("Create Purchase Invoice")
    purchaseinvoice.maskingMoney()
    purchaseinvoice.switchButton()
    purchaseinvoice.getDateNow()
    purchaseinvoice.getDataSupplier()
    purchaseinvoice.processsPlusPlus()
    purchaseinvoice.departmentDropdown()
    //   purchaseinvoice.getDataPurchaseInvoiceStatusPI(function() {
    //         purchaseinvoice.renderGrid()
    //         // console.log("masuk kemare")
    //    })
    purchaseinvoice.names("Non Inventory")
    $(".noninv").show()
    $(".invhide").hide()
    $(".item").prop('disabled', false)
    purchaseinvoice.GetDataPurchaseInvoice()
    purchaseinvoice.GetDataPurchaseInvoiceInventory()
    $("#noninv").css({
        "background-color": "#ffffff",
        "color": "#3ea49d",
        "border-bottom": "4px solid #f4222d"
    })
    // $("#my-checkbox").bootstrapSwitch('toggleDisabled',true,true)

   

}

$(function() {

    purchaseinvoice.init()
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null ||type != null) {
       if (type =="POINV"){
        purchaseinvoice.names("Inventory")
        setTimeout(function(){
            purchaseinvoice.choose('IVENTORY')

        },500)
       }else{
        purchaseinvoice.getDataPurchaseInvoiceStatusPI(function() {
            purchaseinvoice.renderGrid()
        })
       }
    }else{
        purchaseinvoice.getDataPurchaseInvoiceStatusPI(function() {
            purchaseinvoice.renderGrid()
            // console.log("masuk kemare")
        })
    }
})