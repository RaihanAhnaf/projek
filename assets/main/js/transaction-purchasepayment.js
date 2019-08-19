var purchasepayment = {}

purchasepayment.date = ko.observable()
purchasepayment.DatePageBar = ko.observable()
purchasepayment.selectAll = ko.observable(false)
purchasepayment.dataDropDownSupplier = ko.observableArray([])
purchasepayment.dataMasterPurchaseOrder = ko.observableArray([])
purchasepayment.dataMasterPurchasePayment = ko.observableArray([])
purchasepayment.dataMasterPurchasePaymentOriginal = ko.observableArray([])
purchasepayment.dataDropDownSupplierFilter = ko.observableArray([])
purchasepayment.paidBool = ko.observable(false)
purchasepayment.textSearch = ko.observable()
purchasepayment.textSupplierSearch = ko.observable()
purchasepayment.filterindicator = ko.observable()
purchasepayment.setDocumentNumber = ko.observable()
purchasepayment.text = ko.observable()
purchasepayment.backToList = ko.observable(false)
purchasepayment.printPDFListView = ko.observable(false)
purchasepayment.valueDepartment = ko.observable("")
purchasepayment.dataDepartment = ko.observableArray([])
purchasepayment.showAttachment = ko.observable(true)
purchasepayment.names = ko.observable()
var processs = 0;
purchasepayment.payment = ko.observableArray([])
// purchasepayment.payment = [{
//     value: 1110,
//     text: "KAS"
// }, {
//     value: 1120,
//     text: "BANK USD"
// }, {
//     value: 1121,
//     text: "BANK IDR"
// }]
purchasepayment.getDataPaymentAccount = function () {
    model.Processing(true)
    ajaxPost("/transaction/getdatapaymentaccount", {}, function (res) {
        model.Processing(false)
        var data = []
        for (i in res.Data) {
            data.push({
                value: res.Data[i].ACC_Code,
                text: res.Data[i].Account_Name
            })
        }
        purchasepayment.payment(data)
    })
}
purchasepayment.newRecord = function () {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumber: "",
        SupplierCode: "",
        SupplierName: "",
        PaymentAccount: 0,
        PaymentName: "",
        IsInventory: false,
        ListDetail: [],
        Department: "",
        Attachment: "",
        PICMV: false,
        User: ''
    }

    page.ListDetail.push(purchasepayment.listDetail({}))

    return page
}

purchasepayment.listDetail = function (data) {
    var dataTmp = {}

    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.DatePayment = data.DatePayment == undefined ? '' : data.DatePayment
    dataTmp.PoNumber = data.PoNumber == undefined ? '' : data.PoNumber
    dataTmp.Amount = data.Amount == undefined ? 0 : data.Amount
    dataTmp.AlreadyPaid = data.AlreadyPaid == undefined ? 0 : data.AlreadyPaid
    dataTmp.Payment = data.Payment == undefined ? 0 : data.Payment
    dataTmp.Balance = data.Balance == undefined ? 0 : data.Balance
    dataTmp.Pay = data.Pay == undefined ? false : data.Pay

    data = ko.mapping.fromJS(dataTmp)

    data.Payment.subscribe(function (e) {
        var bl = FormatCurrency(data.Amount()) - (FormatCurrency(data.AlreadyPaid()) + FormatCurrency(e))
        if (bl < 0) {
            return swal({
                title: "Warning!",
                text: "your payments are too much",
                type: "info",
                confirmButtonColor: "#3da09a"
            }, function () {
                data.Payment(0)
            });
        } else {
            data.Balance(ChangeToRupiah(bl))
        }
    })

    return data
}

purchasepayment.record = ko.mapping.fromJS(purchasepayment.newRecord())

purchasepayment.getDocumentNumber = function () {
    var res = {
        Number: 1
    }
    var tgl = moment(Date()).format('/DDMMYY/')
    var zr
    if (res.Number < 10) {
        zr = '00'
    } else if (res.Number >= 10 && res.Number < 100) {
        zr = '0'
    } else {
        zr = ''
    }

    purchasepayment.record.DocumentNumber('PP' + tgl + zr + res.Number);
    purchasepayment.setDocumentNumber('PP' + tgl + zr + res.Number)
}

purchasepayment.getDataPurchaseOrder = function (callback) {
    purchasepayment.dataMasterPurchaseOrder([])
    var siki = Date();
    var now = new Date();
    var past = now.setMonth(now.getMonth() - 1, 1);
    var param = {
        DateStart: moment(past).format('YYYY-MM-DD'),
        DateEnd: moment(siki).format('YYYY-MM-DD'),
        IsInventory: purchasepayment.names() == "Inventory"
    }
    ajaxPost('/transaction/getdatapoforpp', param, function (res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchasepayment.dataMasterPurchaseOrder(res.Data)
        model.Processing(false)
        purchasepayment.processsPlusPlus()
        if (typeof callback == "function") callback();
    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

purchasepayment.getDataPurchaseInventory = function (callback) {
    purchasepayment.dataMasterPurchaseOrder([])
    var siki = Date();
    var now = new Date();
    var past = now.setMonth(now.getMonth() - 1, 1);
    var param = {
        DateStart: moment(past).format('YYYY-MM-DD'),
        DateEnd: moment(siki).format('YYYY-MM-DD'),
    }
    ajaxPost('/transaction/getdatapoinventoryforpp', param, function (res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        purchasepayment.dataMasterPurchaseOrder(res.Data)
        model.Processing(false)
        purchasepayment.processsPlusPlus()
        if (typeof callback == "function") callback();
    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

purchasepayment.getAllDataPurchasePayment = function (callback) {
    var param = {}
    if (purchasepayment.filterindicator() == true) {
        param = {
            DateStart: $("#dateStart").data("kendoDatePicker").value(),
            DateEnd: $("#dateEnd").data("kendoDatePicker").value(),
            Filter: true,
            TextSearch: purchasepayment.textSearch(),
            SupplierCode: purchasepayment.textSupplierSearch(),
        }
    } else {
        param = {
            Filter: false,
            DateStart: $("#dateStart").data("kendoDatePicker").value(),
            DateEnd: $("#dateEnd").data("kendoDatePicker").value(),
        }
    }
    model.Processing(true)
    ajaxPost('/transaction/getalldatapurchasepayment', param, function (res) {
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
        purchasepayment.dataMasterPurchasePayment(res.Data)
        purchasepayment.dataMasterPurchasePaymentOriginal(res.Data)
        if (typeof callback == "function") callback();
        purchasepayment.processsPlusPlus()
    })
}

purchasepayment.getSuppplierCode = function () {
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
        purchasepayment.dataDropDownSupplier(res.Data)
        var DataSupplier = res.Data
        for (i in DataSupplier) {
            DataSupplier[i].Kode = DataSupplier[i].Kode + ""
            DataSupplier[i].Name = DataSupplier[i].Kode + " - " + DataSupplier[i].Name
        }
        purchasepayment.dataDropDownSupplierFilter(DataSupplier)
        model.Processing(false)
        purchasepayment.processsPlusPlus()
    })
    // if (model.Processing()== false){

    // }
}

purchasepayment.onChangeSupplierCode = function (value) {
    var data = _.filter(purchasepayment.dataMasterPurchaseOrder(), {
        'SupplierCode': value,
        'Status': 'PI'
    });
    // console.log("-------------", data)
    if (data.length == 0) {
        purchasepayment.reset()
        return swal({
            title: "Warning!",
            text: "There is no po with that supplier code",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    } else {
        purchasepayment.record.SupplierCode(data[0].SupplierCode)
        purchasepayment.record.SupplierName(data[0].SupplierName)
        purchasepayment.record.ListDetail([])
        for (var i = 0; i < data.length; i++) {
            purchasepayment.record.ListDetail.push(purchasepayment.listDetail({}))
            purchasepayment.record.ListDetail()[i].DatePayment(data[i].DateStr)
            purchasepayment.record.ListDetail()[i].PoNumber(data[i].DocumentNumber)
            purchasepayment.record.ListDetail()[i].Amount(data[i].GrandTotalIDR)
            purchasepayment.record.ListDetail()[i].AlreadyPaid(data[i].AlreadyPaid)
            purchasepayment.record.ListDetail()[i].Balance(data[i].GrandTotalIDR - data[i].AlreadyPaid)
        }
        $("[name='pay']").bootstrapSwitch();
        purchasepayment.maskingMoney()
    }
}

purchasepayment.reloadGrid = function () {
    var name = purchasepayment.names();
    console.log(name);
    var callback = function () {
        purchasepayment.renderGrid();
    }
    if (name == "Inventory") {
        purchasepayment.getAllDataPurchasePayment(callback);
    } else {
        purchasepayment.getDataPurchaseInventory(callback);
    }
}

purchasepayment.renderGrid = function () {
    var mydata = name == "Non Inventory" ? purchasepayment.dataMasterPurchasePayment() : purchasepayment.dataMasterPurchaseOrder();
    var name = purchasepayment.names()
    var dsGrid = [];
    for (var i = 0; i < mydata.length; i++) {
        var data = mydata[i];
        if (data.IsInventory && name == "Inventory") dsGrid.push(data);
        if (!data.IsInventory && name == "Non Inventory") dsGrid.push(data);
    }
    dsGrid = mydata;

    var columns = [{
        title: 'Action',
        width: 55,
        template: "<button onclick='purchasepayment.viewDetail(\"#: _id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>"
    },
    {
        field: 'DateStr',
        title: 'Date',
        width: 160,
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 160,
    }, {
        field: 'SupplierName',
        title: 'Supplier Name',
        width: 200,
    },
    {
        title: 'Amount',
        width: 170,
        attributes: {
            "class": "rightAction",
        },
        template: function (e) {
            var list = e.ListDetail
            var value = _.sumBy(list, function (o) { return o.Amount; });
            return kendo.toString(value, 'n2')
        }
    },
    {
        title: 'Payment',
        width: 170,
        attributes: {
            "class": "rightAction",
        },
        template: function (e) {
            var list = e.ListDetail
            var value = _.sumBy(list, function (o) { return o.Payment; });
            return kendo.toString(value, 'n2')
        }
    },
    {
        field: 'BalanceAll',
        title: 'Balance',
        width: 170,
        template: "#=ChangeToRupiah(BalanceAll)#",
        attributes: {
            "class": "rightAction",
        },
    }
    ]

    $('#gridListPurchasePayment').kendoGrid({
        dataSource: {
            data: dsGrid,
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
            ProActive.kendoExcelRender(e, "PurchasePayment", function (row, sheet) {
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "data") {
                        if (ci == 3 || ci == 4 || ci == 5) {
                            cell.format = "#,##0.00_);(#,##0.00);0.00;";
                            cell.hAlign = "right";
                        }
                    }
                }
            });
        },
    })
}

purchasepayment.exportExcel = function () {
    $("#gridListPurchasePayment").getKendoGrid().saveAsExcel();
}
purchasepayment.maskingMoney = function () {
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

purchasepayment.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    purchasepayment.DatePageBar(page)
}

purchasepayment.selectAll = function () {
    $("form input[type='submit']").click(function (e) {
        e.preventDefault();

        // This gets the value of the currently selected option
        var value = $(this).children(":selected").val();

        $(".Checkbox").each(function () {
            if ($(this).is(':checked')) {
                $(this).fadeOut();
                // Do stuff with checked box
            } else {
                // Checkbox isn't checked
            }
        })
    });
}
purchasepayment.search = function (e) {
    // purchasepayment.textSearch(e)
    // purchasepayment.filterindicator(true)
    // purchasepayment.getAllDataPurchasePayment(function() {
    //     purchasepayment.renderGrid()
    // })

    if (purchasepayment.names() == "Non Inventory") {
        purchasepayment.textSearch(e)
        purchasepayment.filterindicator(true)
        purchasepayment.getAllDataPurchasePayment(function () {
            purchasepayment.renderGrid()
        })
    } else {
        purchasepayment.textSearch(e)
        purchasepayment.filterindicator(true)
        purchasepayment.getDataPurchaseInventory(function () {
            purchasepayment.renderGrid()
        })
    }
}

purchasepayment.createdForm = function () {
    purchasepayment.showAttachment(true)
    purchasepayment.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    // purchasepayment.record.DocumentNumber(purchasepayment.setDocumentNumber())
    purchasepayment.reset()
    var dropDown1 = $("#payment").data("kendoDropDownList");
    dropDown1.enable(true);
    var dropDown2 = $("#supliercode").data("kendoDropDownList");
    dropDown2.enable(true);
    purchasepayment.text("Create Purchase Payment")
    $('.space').show()
    purchasepayment.backToList(false)
    purchasepayment.printPDFListView(false)

}

purchasepayment.reset = function () {
    ko.mapping.fromJS(purchasepayment.newRecord(), purchasepayment.record)
    // purchasepayment.getDocumentNumber()
    $(".formInput").val("")
    purchasepayment.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    $('#payment').data('kendoDropDownList').value(-1);
    $('#supliercode').data('kendoDropDownList').value(-1);
    $('#departmenDropdown').data('kendoDropDownList').value(-1);
    $("[name='pay']").bootstrapSwitch();
}

purchasepayment.saveData = function () {
    var change = ko.mapping.toJS(purchasepayment.record)
    purchasepayment.valueDepartment(change.Department)
    for (var i = 0; i < change.ListDetail.length; i++) {
        change.ListDetail[i].DatePayment = purchasepayment.dataMasterPurchaseOrder()[i].DatePosting
    }
    var pay = _.filter(change.ListDetail, ['Pay', true]);
    for (var i = 0; i < pay.length; i++) {
        pay[i].Balance = FormatCurrency(pay[i].Balance)
        pay[i].Payment = FormatCurrency(pay[i].Payment)
        pay[i].AlreadyPaid = FormatCurrency(pay[i].AlreadyPaid)
    }
    // console.log(pay, change)
    if (purchasepayment.valueDepartment() == "") {
        return swal({
            title: "Warning!",
            text: "None Department has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (pay.length == 0 || change.SupplierCode == "") {
        return swal({
            title: "Warning!",
            text: "you didn't pay anything",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (change.PaymentName == "") {
        return swal({
            title: "Warning!",
            text: "you haven't choose the payment ",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    change.ListDetail = pay
    change.PaymentAccount = parseInt(purchasepayment.record.PaymentAccount())
    var c = _.result(_.find(purchasepayment.payment(), function (obj) {
        return obj.value === parseInt(purchasepayment.record.PaymentAccount());
    }), 'text');
    change.PaymentName = c
    change.IsInventory = purchasepayment.names() == "Inventory"
    change.User = userinfo.usernameh()
    change.DatePosting = $('#date').data('kendoDatePicker').value().toISOString()
    change.DateStr = moment($('#date').data('kendoDatePicker').value()).format("DD-MMM-YYYY");
    var formData = new FormData()
    formData.append('data', JSON.stringify(change))
    formData.append('department', purchasepayment.valueDepartment())
    formData.append('typepurchase', purchasepayment.names())
    var attachment = document.getElementById('uploadFile');
    if (attachment.files.length == 0) {
        return swal({
            title: "Warning!",
            text: "you haven't choose the attactment ",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (!ProActive.checkUploadFilesize(attachment.files[0].size, true)) return;
    formData.append("fileUpload", attachment.files[0]);
    var url = "/transaction/savepurchasepayment"
    // var param = {
    //     Data: change,
    //     Department: purchasepayment.valueDepartment(),
    // }
    swal({
        title: "Are you sure?",
        text: "You will save this Purchase Payment " + purchasepayment.names(),
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
            $.ajax({
                url: url,
                data: formData,
                contentType: false,
                dataType: "json",
                mimeType: 'multipart/form-data',
                processData: false,
                type: 'POST',
                success: function (data) {
                    if (data.Status == "OK") {
                        setTimeout(function () {
                            swal({
                                title: "Success!",
                                text: "Data has been saved!",
                                type: "success",
                                confirmButtonColor: "#3da09a"
                            }, function () {
                                window.location.assign("/transaction/purchasepayment")
                                $("#btnSave").attr("disabled", "disabled")
                                $("#btnPrint").attr("disabled", "disabled")
                                $("#btnSavePrint").attr("disabled", "disabled")
                            });

                        }, 100)
                        model.Processing(false)
                    }
                }
            });
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
purchasepayment.downloadAttachment = function () {
    var data = ko.toJS(purchasepayment.record)
    var link = moment(data.DatePosting).format('YYYYMM') + '/' + data.Attachment
    var pom = document.createElement('a');
    pom.setAttribute('href', "/res/docs/" + link);
    pom.setAttribute('download', link);
    pom.click();
}
purchasepayment.printToPdf = function () {
    model.Processing(true)
    var change = ko.mapping.toJS(purchasepayment.record)
    purchasepayment.valueDepartment(change.Department)
    for (var i = 0; i < change.ListDetail.length; i++) {
        change.ListDetail[i].DatePayment = purchasepayment.dataMasterPurchaseOrder()[i].DatePosting
    }
    var pay = _.filter(change.ListDetail, ['Pay', true]);
    for (var i = 0; i < pay.length; i++) {
        pay[i].Balance = FormatCurrency(pay[i].Balance)
        pay[i].Payment = FormatCurrency(pay[i].Payment)
        pay[i].AlreadyPaid = FormatCurrency(pay[i].AlreadyPaid)
    }
    if (pay.length == 0 || change.SupplierCode == "") {
        return swal({
            title: "Warning!",
            text: "you didn't pay anything",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (change.PaymentName == "") {
        return swal({
            title: "Warning!",
            text: "you haven't choose the payment ",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    change.ListDetail = pay
    change.PaymentAccount = parseInt(purchasepayment.record.PaymentAccount())

    var c = _.result(_.find(purchasepayment.payment(), function (obj) {
        return obj.value === parseInt(purchasepayment.record.PaymentAccount());
    }), 'text');
    change.PaymentName = c
    change.User = userinfo.usernameh()
    change.DatePosting = $('#date').data('kendoDatePicker').value().toISOString()
    change.DateStr = moment($('#date').data('kendoDatePicker').value()).format("DD-MMM-YYYY");

    var param = {
        Data: change
    }
    ajaxPost("/transaction/exporttopdfpurchasepayment", param, function (e) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/purchasepayment/' + e, '_blank');
        tabOrWindow.focus();
    })
}
purchasepayment.printListToPdf = function (e) {
    model.Processing(true)
    var param = {
        Id: purchasepayment.record.ID(),
    }
    ajaxPost("/transaction/exporttopdflistpurchasepayment", param, function (e) {
        model.Processing(false)
        var taborWindow = window.open('/res/docs/purchasepayment/' + e, '_blank');
        taborWindow.focus();
    })
}

purchasepayment.savePrint = function () {
    var change = ko.mapping.toJS(purchasepayment.record)
    purchasepayment.valueDepartment(change.Department)
    for (var i = 0; i < change.ListDetail.length; i++) {
        change.ListDetail[i].DatePayment = purchasepayment.dataMasterPurchaseOrder()[i].DatePosting
    }
    var pay = _.filter(change.ListDetail, ['Pay', true]);
    for (var i = 0; i < pay.length; i++) {
        pay[i].Balance = FormatCurrency(pay[i].Balance)
        pay[i].Payment = FormatCurrency(pay[i].Payment)
        pay[i].AlreadyPaid = FormatCurrency(pay[i].AlreadyPaid)
    }
    if (purchasepayment.valueDepartment() == "") {
        return swal({
            title: "Warning!",
            text: "None Department has selected",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (pay.length == 0 || change.SupplierCode == "") {
        return swal({
            title: "Warning!",
            text: "you didn't pay anything",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    if (change.PaymentName == "") {
        return swal({
            title: "Warning!",
            text: "you haven't choose the payment ",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    change.ListDetail = pay
    change.PaymentAccount = parseInt(purchasepayment.record.PaymentAccount())

    var c = _.result(_.find(purchasepayment.payment(), function (obj) {
        return obj.value === parseInt(purchasepayment.record.PaymentAccount());
    }), 'text');
    change.PaymentName = c
    change.User = userinfo.usernameh()
    change.DatePosting = $('#date').data('kendoDatePicker').value().toISOString()
    change.DateStr = moment($('#date').data('kendoDatePicker').value()).format("DD-MMM-YYYY");
    var formData = new FormData()
    formData.append('data', JSON.stringify(change))
    formData.append('department', purchasepayment.valueDepartment())
    var attachment = document.getElementById('uploadFile');
    if (attachment.files.length == 0) {
        if (attachment[i].files[0])
            if (!ProActive.checkUploadFilesize(attachment[i].files[0].size,
                "Attachment file size cannot exceed $size (Row #" + (i + 1) + ")!")) return;
        return swal({
            title: "Warning!",
            text: "you haven't choose the attactment ",
            type: "info",
            confirmButtonColor: "#3da09a"
        })
    }
    formData.append("fileUpload", attachment.files[0]);
    var url = "/transaction/savepurchasepayment"
    // var param = {
    //     Data: change,
    //     Department: purchasepayment.valueDepartment(),
    // }
    swal({
        title: "Are you sure?",
        text: "You will save this Purchase Payment " + purchasepayment.names(),
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
            model.Processing(true)
            $.ajax({
                url: url,
                data: formData,
                contentType: false,
                dataType: "json",
                mimeType: 'multipart/form-data',
                processData: false,
                type: 'POST',
                success: function (data) {
                    if (data.Status == "OK") {
                        setTimeout(function () {
                            swal({
                                title: "Success!",
                                text: "Data has been saved!",
                                type: "success",
                                confirmButtonColor: "#3da09a"
                            }, function (e) {
                                purchasepayment.printToPdf()
                                setTimeout(function () {
                                    window.location.assign("/transaction/purchasepayment")
                                }, 200)
                            });

                        }, 100)
                        model.Processing(false)
                    }
                }
            });
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

purchasepayment.onChangePaid = function (n, value) {
    var index = n.replace("pay", "");
    purchasepayment.record.ListDetail()[index].Pay(value)
    if (value) {
        $('#payment' + index).removeAttr("disabled");
    } else {
        $('#payment' + index).attr('disabled', 'disabled')
    }
}
purchasepayment.onChangePayment = function (e) {
    var e = parseInt(e)
    var accName = _.filter(purchasepayment.payment(), ['value', e])[0].text
    purchasepayment.record.PaymentName(accName)
}

purchasepayment.checkRole = function () {
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#date').data('kendoDatePicker').enable(condition);
}

purchasepayment.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}

purchasepayment.viewDetail = function (e) {
    // var allData = purchasepayment.dataMasterPurchasePaymentOriginal()
    var allData = purchasepayment.gridDataSource()
    var data = _.find(allData, function (o) {
        return o._id == e;
    });
    if (data.length == 0) {
        return swal({
            title: "Error!",
            text: "Data is not found",
            type: "error",
            confirmButtonColor: "#3da09a"
        });
    }
    // console.log(data)
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    ko.mapping.fromJS(data, purchasepayment.record)
    purchasepayment.record.ID(data._id)
    _.each(purchasepayment.record.ListDetail(), function (v, i) {
        purchasepayment.record.ListDetail()[i].AlreadyPaid(v.Amount() - (v.Payment() + v.Balance()))
        var n = moment(v.DatePayment()).format('DD-MMM-YYYY')
        purchasepayment.record.ListDetail()[i].DatePayment(n)
    });
    // $('#departmenDropdown').data("kendoDropdownList").value(data.Department);
    $("#ListPurchasePayment").removeClass("active")
    $("#CreatePurchasePayment").addClass("active")
    var dropDown1 = $("#payment").data("kendoDropDownList");
    dropDown1.enable(false);
    var dropDown2 = $("#supliercode").data("kendoDropDownList");
    dropDown2.enable(false);
    var dropDown3 = $("#departmenDropdown").data("kendoDropDownList");
    dropDown3.enable(false);
    purchasepayment.text("View Purchase Payment")
    purchasepayment.backToList(true)
    $('.space').hide()
    $("[name='pay']").bootstrapSwitch({
        disabled: true,
        state: true
    });
    $('.nav-tabs a[href="#CreatePurchasePayment"]').tab('show');
    purchasepayment.showAttachment(false)
    purchasepayment.maskingMoney()
    purchasepayment.printPDFListView(true)
}

purchasepayment.backList = function () {
    purchasepayment.showAttachment(true)
    $('.nav-tabs a[href="#ListPurchasePayment"]').tab('show')
    $("#ListPurchasePayment").addClass("active");
    $("#CreatePurchasePayment").removeClass("active");
}
purchasepayment.departmentDropdown = function () {
    $("#departmenDropdown").html("")
    ajaxPost("/transaction/getdatadepartment", {}, function (res) {
        purchasepayment.dataDepartment(res.Data)
        // $("#departmenDropdown").kendoDropDownList({
        //     value: purchasepayment.valueDepartment(),
        //     filter: "contains",
        //     dataTextField: "DepartmentName",
        //     dataValueField: "DepartmentName",
        //     dataSource: res.Data,
        //     optionLabel:'Select one',
        //     change:function(e){
        //         var dataitem = this.dataItem();
        //         purchasepayment.valueDepartment(dataitem.DepartmentName)
        //     }   
        // });
    })
}

purchasepayment.choose = function (choosenote) {
    if (choosenote == "ACTIVA") {
        purchasepayment.names("Non Inventory")
        purchasepayment.getDataPurchaseOrder()

        $('.invhidepi').hide()
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


    } else {
        purchasepayment.names("Inventory")
        $('.invhidepi').show()
        purchasepayment.getDataPurchaseInventory()
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
    purchasepayment.reset()
    $('.nav-tabs a[href="#ListPurchasePayment"]').tab('show')
}


purchasepayment.init = function () {
    var now = new Date()
    purchasepayment.getDataPaymentAccount()
    ProActive.KendoDatePickerRange();
    purchasepayment.getSuppplierCode()
    purchasepayment.getDataPurchaseOrder()
    purchasepayment.departmentDropdown()
    purchasepayment.setDate()
    purchasepayment.getDateNow()
    purchasepayment.date(moment(now).format("DD MMMM YYYY"))
    purchasepayment.checkRole()
    purchasepayment.processsPlusPlus()
    $("#noninv").css({
        "background-color": "#ffffff",
        "color": "#3ea49d",
        "border-bottom": "4px solid #f4222d"
    });
    purchasepayment.names("Non Inventory")
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var type = url.searchParams.get("type")
    if (num != null) {
        if (type == "POINV") {
            purchasepayment.choose('IVENTORY')
        }
    }
    purchasepayment.search(true);
    
    ProActive.GlobalSearch("gridListPurchasePayment",
        [
            "DocumentNumber",
            "SupplierCode",
            "SupplierName",
            "Remark"
        ]
    );
}
purchasepayment.createPPFromPOInvSummary = function () {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var stat = url.searchParams.get("stat");
    var type = url.searchParams.get("type")
    if (num != null) {
        var data = _.filter(purchasepayment.dataMasterPurchasePayment(), {
            'SupplierCode': num,
        });
        if (stat == "PAID") {
            purchasepayment.viewDetail(data[0]._id)
        } else {
            purchasepayment.createdForm()
            $('.nav-tabs a[href="#CreatePurchasePayment"]').tab('show');
            $("#supliercode").data("kendoDropDownList").value(num)
            purchasepayment.onChangeSupplierCode(num)
        }
    }
}
purchasepayment.processsPlusPlus = function () {
    processs += 1
    // console.log(processs)
    if (processs >= 4) {
        setTimeout(function () {
            purchasepayment.createPPFromPOInvSummary()
        }, 300)
        // setTimeut(purchasepayment.createPPFromPOInvSummary(),10000) 
    }
}
$(function () {
    $('.invhidepi').hide()
    $(document).on('change', ':file', function () {
        var input = $(this),
            numFiles = input.get(0).files ? input.get(0).files.length : 1,
            label = input.val().replace(/\\/g, '/').replace(/.*\//, '');
        input.trigger('fileselect', [numFiles, label]);
    });
    $(':file').on('fileselect', function (event, numFiles, label) {

        var input = $(this).parents('.input-group').find(':text'),
            log = numFiles > 1 ? numFiles + ' files selected' : label;

        if (input.length) {
            input.val(log);
        } else {
            if (log) alert(log);
        }

    });
    purchasepayment.text("Create Purchase Payment")
    purchasepayment.init()
    $('#tablePP').on('switchChange.bootstrapSwitch', 'input[name="pay"]', function (event, state) {
        purchasepayment.onChangePaid(event.currentTarget.id, state)
    });

})


// ==== New Implementation ====
// Below functions overwrites any functions above with same name
purchasepayment.gridDataSource = ko.observableArray([]);
purchasepayment.reloadDataSource = function (callback, fromSearch) {
    var mode = purchasepayment.names();
    fromSearch = fromSearch === true;
    var startdate = $('#dateStart').data('kendoDatePicker').value();
    var enddate = $('#dateEnd').data('kendoDatePicker').value();
    // var startdate = purchasepayment.DateStart();
    // var enddate = purchasepayment.DateEnd();

    var param = {
        DateStart: moment(startdate).format('YYYY-MM-DD'),
        DateEnd: moment(enddate).format('YYYY-MM-DD'),
        SupplierCode: purchasepayment.textSupplierSearch(),
        Filter: true,
    }

    if (mode == "Inventory") {
        param.IsInventory = true
    } else{
        param.IsInventory = false
    }

    // var url = mode == "Inventory" ? '/purchase/getdatapoinventoryforpp' : '/purchase/getdatapoforpp';
    var url = '/transaction/getdatapurchasepayment'
    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            });
            return
        }

        // for (i in res.Data) {
        //     var data = res.Data[i];
        //     data.BalanceAll = data.GrandTotalIDR - data.AlreadyPaid;
        //     data.BalanceAll = data.Amount - data.AlreadyPaid;
        // }

        purchasepayment.gridDataSource(res.Data);

        if (typeof callback == "function") callback();
    }, function () {
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}
purchasepayment.renderGrid = function () {
    var grid = $('#gridListPurchasePayment');
    if (typeof grid.data('kendoGrid') !== 'undefined') {
        grid.data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: purchasepayment.gridDataSource(),
        }))
        return
    }

    var columns = [{
        title: 'Action',
        width: 80,
        template: "<button onclick='purchasepayment.viewDetail(\"#: _id #\")' class='btn btn-sm btn-default'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>"
    },
    {
        field: 'DateStr',
        title: 'Date',
        width: 160,
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 160,
    }, {
        field: 'SupplierName',
        title: 'Supplier Name',
        width: 200,
    },
    {
        title: 'Amount',
        width: 170,
        attributes: {
            "class": "rightAction",
        },
        template: function (e) {
            var list = e.ListDetail
            // var value = _.sumBy(list, function (o) { return o.AmountIDR; }) || 0;
            var value = _.sumBy(list, function (o) { return o.Amount; }) || 0;
            // value += e.VAT;
            return kendo.toString(value, 'n2')
        }
    },
    {
        title: 'Payment',
        width: 170,
        attributes: {
            "class": "rightAction",
        },
        template: function (e) {
            var list = e.ListDetail
            // var value = _.sumBy(list, function (o) { return o.PaymentIDR; }) || 0;
            var value = _.sumBy(list, function (o) { return o.Payment; }) || 0;
            return kendo.toString(value, 'n2')
        }
    },
    {
        field: 'BalanceAll',
        title: 'Balance',
        width: 170,
        template: "#=ChangeToRupiah(BalanceAll)#",
        attributes: {
            "class": "rightAction",
        },
    }
    ]

    grid.kendoGrid({
        dataSource: {
            data: purchasepayment.gridDataSource(),
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
            ProActive.kendoExcelRender(e, "PurchasePayment", function (row, sheet) {
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "data") {
                        if (ci == 3 || ci == 4 || ci == 5) {
                            cell.format = "#,##0.00_);(#,##0.00);0.00;";
                            cell.hAlign = "right";
                        }
                    }
                }
            });
        },
    });
}
purchasepayment.search = function (notFromSearchButton) {
    purchasepayment.filterindicator(true);
    model.Processing(true);
    purchasepayment.reloadDataSource(function() { 
        purchasepayment.renderGrid();
        model.Processing(false); 
    }, notFromSearchButton === true);
}