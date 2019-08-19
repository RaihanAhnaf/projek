var salespayment = {}
var processs = 0;
salespayment.DatePageBar = ko.observable()
salespayment.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
salespayment.DateEnd = ko.observable(new Date)
salespayment.text = ko.observable()
salespayment.setDocumentNumber = ko.observable()
salespayment.selectAll = ko.observable()
salespayment.sequenceNumber = ko.observable()
salespayment.filterindicator = ko.observable()
salespayment.textSearch = ko.observable()
salespayment.save = ko.observable(false)
salespayment.print = ko.observable(false)
salespayment.printPDFListView = ko.observable(false)
salespayment.saveAndPrint = ko.observable(false)
salespayment.backToList = ko.observable(false)
salespayment.dataMasterCustomer = ko.observableArray([])
salespayment.dataDropDownCustomer = ko.observableArray([])
salespayment.dataMasterSalesPayment = ko.observableArray([])
salespayment.dataMasterSalesPaymentOriginal = ko.observableArray([])
salespayment.dataDropDownCustomerFilter = ko.observableArray([])
salespayment.textCustomerSearch = ko.observable()
salespayment.showAttachment = ko.observable(true)
salespayment.filterStatus = ko.observableArray([])
// salespayment.filterStatus = [{
//     value: 1110,
//     text: "KAS"
// }, {
//     value: 1120,
//     text: "BANK USD"
// }, {
//     value: 1121,
//     text: "BANK IDR"
// }]
salespayment.getDataPaymentAccount = function(){
    model.Processing(true)
    ajaxPost("/transaction/getdatapaymentaccount", {}, function(res){
        model.Processing(false)
        var data = []
        for(i in res.Data){
            data.push({
                value: res.Data[i].ACC_Code,
                text: res.Data[i].Account_Name
            })
        }
        salespayment.filterStatus(data)
    })
}
salespayment.newRecord = function () {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumber: "",
        CustomerCode: "",
        CustomerName: "",
        PaymentAccount: "",
        PaymentName: "",
        Attachment:"",
        ListDetail: []
    }

    page.ListDetail.push(salespayment.listDetail({}))

    return page
}

salespayment.listDetail = function (data) {
    var dataTmp = {}

    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.DatePayment = data.DatePayment == undefined ? '' : data.DatePayment
    dataTmp.InvNumber = data.InvNumber == undefined ? '' : data.InvNumber
    dataTmp.Amount = data.Amount == undefined ? 0 : data.Amount
    dataTmp.Receive = data.Receive == undefined ? 0 : data.Receive
    dataTmp.Balance = data.Balance == undefined ? 0 : data.Balance
    dataTmp.AlreadyPaid = data.AlreadyPaid == undefined ? 0 : data.AlreadyPaid
    dataTmp.Pay = data.Pay == undefined ? false : data.Pay

    var data = ko.mapping.fromJS(dataTmp)

    data.Receive.subscribe(function (e) {
        var bl = FormatCurrency(data.Amount()) - (FormatCurrency(data.AlreadyPaid()) + FormatCurrency(e))
        if (bl < 0) {
            return swal({
                title: "Warning!",
                text: "your payments are too match",
                type: "info",
                confirmButtonColor: "#3da09a"
            }, function () {
                data.Receive(0)
            });
        } else {
            data.Balance(ChangeToRupiah(bl))
        }
    })

    return data
}

salespayment.record = ko.mapping.fromJS(salespayment.newRecord())


salespayment.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    salespayment.DatePageBar(page)
}

salespayment.getDocumentNumber = function () {
    model.Processing(true)
    ajaxPost('/transaction/getlastnumbersp', {}, function (res) {
        var tgl = moment(Date()).format('/DDMMYY/')
        var zr
        if (res.Number < 10) {
            zr = '00'
        } else if (res.Number >= 10 && res.Number < 100) {
            zr = '0'
        } else {
            zr = ''
        }

        salespayment.record.DocumentNumber('SP' + tgl + zr + res.Number);
        salespayment.sequenceNumber(res.Number)
        salespayment.setDocumentNumber('SP' + tgl + zr + res.Number);
        salespayment.processsPlusPlus()
    })
}

salespayment.getDataInvoice = function () {
    var now = Date();
    var newNow = new Date();
    var past = newNow.setMonth(newNow.getMonth() - 1, 1);
    var param = {
        DateStart: moment(past).format('YYYY-MM-DD'),
        DateEnd: moment(now).format('YYYY-MM-DD'),
    }
    ajaxPost('/transaction/getdatainvoiceforsp', param, function (res) {
        if (res.IsError === "true") {
            swal({
                title: "Error!", 
                text: res.Message, 
                type: "error",
                confirmButtonColor: "#3da09a"})
            return
        }
        salespayment.dataMasterCustomer(res.Data)
        model.Processing(false)
        salespayment.processsPlusPlus()
    }, function () {
        swal({
            title: "Error!", 
            text: "Unknown error, please try again", 
            type: "error",
            confirmButtonColor: "#3da09a"})
    })
}

salespayment.getAllDataSalesPayment = function (callback) {
    var param = {}
    if (salespayment.filterindicator() == true) {
        param = {
            DateStart: $("#dateStart").data("kendoDatePicker").value(),
            DateEnd: $("#dateEnd").data("kendoDatePicker").value(),
            Filter: true,
            TextSearch: salespayment.textSearch(),
            CustomerCode: salespayment.textCustomerSearch(),
        }
    } else {
        param = {
            DateStart: $("#dateStart").data("kendoDatePicker").value(),
            DateEnd: $("#dateEnd").data("kendoDatePicker").value(),
            Filter: false,
        }
    }
    model.Processing(true)
    ajaxPost('/transaction/getalldatasalespayment', param, function (res) {
        if (res.IsError) {
            swal({
                title: "No Data Found!", 
                text: res.Message, 
                type: "warning",
                confirmButtonColor: "#3da09a"})
            $('#textSearch').val("")
            res.Data = []
            // return
        }
        salespayment.processsPlusPlus()
        model.Processing(false)
        salespayment.dataMasterSalesPayment(res.Data)
        salespayment.dataMasterSalesPaymentOriginal(res.Data)
        callback()
    })
}

salespayment.getCustomerCode = function () {
    model.Processing(true)

    ajaxPost('/transaction/getcustomer', {}, function (res) {

        if (res.Total === 0) {
            swal({
                title: "Error!", 
                text: res.Message, 
                type: "error",
                confirmButtonColor: "#3da09a"})
            return
        }
        salespayment.processsPlusPlus()
        salespayment.dataDropDownCustomer(res.Data)
        var DataCustomer = res.Data
        for (i in DataCustomer) {
            DataCustomer[i].Kode = DataCustomer[i].Kode + ""
            DataCustomer[i].Name = DataCustomer[i].Kode + "-" + DataCustomer[i].Name
        }
        salespayment.dataDropDownCustomerFilter(DataCustomer)
        model.Processing(false)
    })
}

salespayment.renderGrid = function () {
    var mydata = salespayment.dataMasterSalesPayment()
    var columns = [
    {
        title: 'View',
        width: 100,
        template: "<button onclick='salespayment.viewList(\"#: _id #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>"
    },
    {
        field: 'DateStr',
        title: 'Date',
        width: 200,
    },{
        field: 'StoreLocationName',
        title: 'Location',
        width: 160,
        template : function(dt){
            return dt.StoreLocationName + " (" + dt.StoreLocationId + ")"
        },
    },{
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 200,
    }, {
        field: 'CustomerName',
        title: 'Customer Name',
        width: 200,
    }, 
    {
        title: 'Amount',
        width: 170,
        attributes: {
            "class": "rightAction",
        },
        template : function(e){
            var list = e.ListDetail
            var value = _.sumBy(list, function(o) { return o.Amount; });
            return kendo.toString(value, 'n2')
        }
    },
    {
        title: 'Receive',
        width: 170,
        attributes: {
            "class": "rightAction",
        },
        template : function(e){
            var list = e.ListDetail
            var value = _.sumBy(list, function(o) { return o.Receive; });
            return kendo.toString(value, 'n2')
        }
    },{
        field: 'BalanceAll',
        title: 'Balance',
        width: 200,
        template: "#=ChangeToRupiah(BalanceAll)#",
         attributes: {
            "class": "rightAction",
        },
    }]

    $('#gridListSalesPayment').kendoGrid({
        dataSource: {
            data: mydata,
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
            ProActive.kendoExcelRender(e, "SalesCreditPayment", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
                        if (ci == 0) {
                            cell.format = "dd-MM-yyyy"; 
                        }
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
salespayment.exportExcel = function () {
    $("#gridListSalesPayment").getKendoGrid().saveAsExcel();
}

salespayment.search = function (e) {
    salespayment.filterindicator(true)
    salespayment.getAllDataSalesPayment(function () {
        salespayment.renderGrid()
    })
}

salespayment.createdForm = function () {
    salespayment.showAttachment(true)
    // salespayment.record.DocumentNumber(salespayment.setDocumentNumber())
    salespayment.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    salespayment.resetView()
    salespayment.save(true)
    salespayment.print(true)
    salespayment.saveAndPrint(true)
    salespayment.backToList(false)
    salespayment.enableView()
    salespayment.printPDFListView(false)
    $("[name='pay']").bootstrapSwitch();

}

salespayment.resetView = function () {
    ko.mapping.fromJS(salespayment.newRecord(), salespayment.record)
    $(".formInput").val("")
    salespayment.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
    salespayment.record.DocumentNumber(salespayment.setDocumentNumber())
    $("#receive").data('kendoDropDownList').value(-1);
    $("#customercode").data("kendoDropDownList").value(-1);
    $("[name='pay']").bootstrapSwitch();
}

salespayment.disableView = function () {
    $(".formInput").attr("disabled", "disabled")
    $('#date').data('kendoDatePicker').enable(false);
    var dropDown1 = $("#customercode").data("kendoDropDownList");
    dropDown1.enable(false);
    var dropDown2 = $("#receive").data("kendoDropDownList");
    dropDown2.enable(false);
    $("[name='pay']").bootstrapSwitch({
        disabled: true,
        state: true
    });
}

salespayment.enableView = function () {
    $(".formInput").attr("disabled", "disabled")
    var condition = (userinfo.rolenameh() == 'supervisor' || userinfo.rolenameh() == 'administrator')
    $('#date').data('kendoDatePicker').enable(condition);
    var dropDown1 = $("#customercode").data("kendoDropDownList");
    dropDown1.enable(true);
    var dropDown2 = $("#receive").data("kendoDropDownList");
    dropDown2.enable(true);
    $("[name='pay']").bootstrapSwitch({
        disabled: false
    });
}

salespayment.backList = function () {
    salespayment.showAttachment(true)
    $('.nav-tabs a[href="#ListSalesPayment"]').tab('show')
    $("#ListSalesPayment").addClass("active");
    $("#CreateSalesPayment").removeClass("active");
}

salespayment.viewList = function (e) {
    salespayment.showAttachment(false)
    var allData = salespayment.dataMasterSalesPaymentOriginal()
    var data = _.find(allData, function (o) {
        return o._id == e;
    });
    if(data.Attachment!=""){
        $('#downloadAttachment').prop('disabled', false);
    }
    if(data.Attachment==""|| data.Attachment== undefined){ 
        $('#downloadAttachment').prop('disabled', true);
    }
    for (var i = 0; i < data.ListDetail.length; i++) {
        data.ListDetail[i].DatePayment = moment(data.ListDetail[i].DatePayment).format('DD-MMM-YYYY')
    }
    data.DatePosting = moment(data.DatePosting).format('DD-MMM-YYYY')
    ko.mapping.fromJS(data, salespayment.record)
    _.each(salespayment.record.ListDetail(), function(v, i) {
        salespayment.record.ListDetail()[i].Receive(0)
        var n = moment(v.DatePayment()).format('DD-MMM-YYYY')
        salespayment.record.ListDetail()[i].DatePayment(n)
    });
    salespayment.record.ID(data._id)
    $("#ListSalesPayment").removeClass("active")
    $("#CreateSalesPayment").addClass("active")
    salespayment.text("View List Sales Payment")
    salespayment.save(false)
    salespayment.print(false)
    salespayment.saveAndPrint(false)
    salespayment.backToList(true)
    salespayment.disableView()
    salespayment.maskingMoney()
    salespayment.printPDFListView(true)
}

salespayment.saveData = function () {
    var change = ko.mapping.toJS(salespayment.record)
    for (var i = 0; i < change.ListDetail.length; i++) {
        change.ListDetail[i].DatePayment = salespayment.dataMasterCustomer()[i].DateCreated
    }
    var pay = _.filter(change.ListDetail, ['Pay', true]);
    for (var i = 0; i < pay.length; i++) {
        pay[i].Balance = FormatCurrency(pay[i].Balance)
        pay[i].Receive = FormatCurrency(pay[i].Receive)
        pay[i].AlreadyPaid = FormatCurrency(pay[i].AlreadyPaid)
    }
    change.ListDetail = pay
    change.PaymentAccount = parseInt(salespayment.record.PaymentAccount())

    var c = _.result(_.find(salespayment.filterStatus(), function (obj) {
        return obj.value === parseInt(salespayment.record.PaymentAccount());
    }), 'text');
    change.PaymentName = c

    change.User = userinfo.usernameh()
    change.DatePosting = $('#date').data('kendoDatePicker').value().toISOString()
    change.DateStr = moment($('#date').data('kendoDatePicker').value()).format("DD-MMM-YYYY");
    // console.log(change)
    if (change.PaymentName == "" || change.PaymentName == undefined) {
        return swal({
            title: 'Warning', 
            text: 'you have not choose the Receiver', 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    if (change.CustomerCode == "") {
        return swal({
            title: 'Warning', 
            text: 'you have not choose the Receiver', 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    if (pay.length == 0) {
        return swal({
            title: 'Warning', 
            text: "you didn't pay anything", 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    var url = "/transaction/savesalespayment"
    var formData = new FormData()
    formData.append('data', JSON.stringify(change))
    formData.append('LastNumber', salespayment.sequenceNumber())
    var attachment = document.getElementById('uploadFile');
    if (attachment.files.length!=0){
        if (attachment.files[0])
        if (!ProActive.checkUploadFilesize(attachment.files[0].size, true)) return;
        formData.append("fileUpload", attachment.files[0]);
    }
    // var param = {
    //     Data: change,
    //     LastNumber: salespayment.sequenceNumber()
    // }
    swal({
        title: "Are you sure?",
        text: "You will save this Sales Payment",
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
                    setTimeout(function () {
                        swal({
                            title: "Success!",
                            text: "Data has been saved!",
                            type: "success",
                            confirmButtonColor: "#3da09a"
                        }, function () {
                            window.location.assign("/transaction/salespayment")
                            $("#btnSave").attr("disabled", "disabled")
                            $("#btnPrint").attr("disabled", "disabled")
                            $("#btnSavePrint").attr("disabled", "disabled")
                        });
                    }, 100)
                    model.Processing(false)
                }
            });
            // ajaxPost(url, param, function (e) {
            //     setTimeout(function () {
            //         swal({
            //             title: "Success!",
            //             text: "Data has been saved!",
            //             type: "success",
            //             confirmButtonColor: "#3da09a"
            //         }, function () {
            //             window.location.assign("/transaction/salespayment")
            //             $("#btnSave").attr("disabled", "disabled")
            //             $("#btnPrint").attr("disabled", "disabled")
            //             $("#btnSavePrint").attr("disabled", "disabled")
            //         });
            //     }, 100)
            //     model.Processing(false)
            // })
        } else {
            swal({
                title: "Cancelled", 
                text: "", 
                type: "error",
                confirmButtonColor: "#3da09a"});
        }
    });
}
salespayment.downloadAttachment = function(){
    var data = ko.toJS(salespayment.record)
    var link = moment(data.DatePosting).format('YYYYMM') + '/' +data.Attachment
    var pom = document.createElement('a');
    pom.setAttribute('href', "/res/docs/" + link);
    pom.setAttribute('download', link);
    pom.click();
}
 salespayment.printToPdf = function () {
     model.Processing(true)
     var change = ko.mapping.toJS(salespayment.record)
    for (var i = 0; i < change.ListDetail.length; i++) {
        var dataMasterList =_.filter(salespayment.dataMasterCustomer(), function(o) { return o.CustomerCode== change.CustomerCode});
        change.ListDetail[i].DatePayment = dataMasterList[i].DateCreated
    }
    var pay = _.filter(change.ListDetail, ['Pay', true]);
    for (var i = 0; i < pay.length; i++) {
        pay[i].Balance = FormatCurrency(pay[i].Balance)
        pay[i].Receive = FormatCurrency(pay[i].Receive)
        pay[i].AlreadyPaid = FormatCurrency(pay[i].AlreadyPaid)
    }
    change.ListDetail = pay
    change.PaymentAccount = parseInt(salespayment.record.PaymentAccount())

    var c = _.result(_.find(salespayment.filterStatus(), function (obj) {
        return obj.value === parseInt(salespayment.record.PaymentAccount());
    }), 'text');
    change.PaymentName = c

    change.User = userinfo.usernameh()
    change.DatePosting = $('#date').data('kendoDatePicker').value().toISOString()
    change.DateStr = moment($('#date').data('kendoDatePicker').value()).format("DD-MMM-YYYY");

    if (change.PaymentName == "" || change.PaymentName == undefined) {
        return swal({
            title: 'Warning', 
            text: 'you have not choose the Receiver', 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    if (change.CustomerCode == "") {
        return swal({
                title: 'Warning', 
                text: 'you have not choose the Receiver', 
                type: 'info',
                confirmButtonColor: "#3da09a"})
    }
    if (pay.length == 0) {
        return swal({
            title: 'Warning', 
            text: "you didn't pay anything", 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
     var param = {
         Data : change
     }
     ajaxPost("/transaction/exporttopdfsalespayment", param, function (e) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/salespayment/' + e, '_blank');
        tabOrWindow.focus();
     })
     return true
 }

 salespayment.printToPdfListView = function(e) {
    model.Processing(true)
     var param = {
         Id: salespayment.record.ID(),
     }
     ajaxPost("/transaction/exporttopdfsalespaymentlistview", param, function(e) {
         model.Processing(false)
         var tabOrWindow = window.open('/res/docs/salespayment/' + e, '_blank');
         tabOrWindow.focus();
     })
 }

 salespayment.SavePrint = function() {
    var change = ko.mapping.toJS(salespayment.record)
    for (var i = 0; i < change.ListDetail.length; i++) {
        change.ListDetail[i].DatePayment = salespayment.dataMasterCustomer()[i].DateCreated
    }
    var pay = _.filter(change.ListDetail, ['Pay', true]);
    for (var i = 0; i < pay.length; i++) {
        pay[i].Balance = FormatCurrency(pay[i].Balance)
        pay[i].Receive = FormatCurrency(pay[i].Receive)
        pay[i].AlreadyPaid = FormatCurrency(pay[i].AlreadyPaid)
    }
    change.ListDetail = pay
    change.PaymentAccount = parseInt(salespayment.record.PaymentAccount())

    var c = _.result(_.find(salespayment.filterStatus(), function (obj) {
        return obj.value === parseInt(salespayment.record.PaymentAccount());
    }), 'text');
    change.PaymentName = c

    change.User = userinfo.usernameh()
    change.DatePosting = $('#date').data('kendoDatePicker').value().toISOString()
    change.DateStr = moment($('#date').data('kendoDatePicker').value()).format("DD-MMM-YYYY");

    if (change.PaymentName == "" || change.PaymentName == undefined) {
        return swal({
            title: 'Warning', 
            text: 'you have not choose the Receiver', 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    if (change.CustomerCode == "") {
        return swal({
            title: 'Warning', 
            text: 'you have not choose the Receiver', 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    if (pay.length == 0) {
        return swal({
            title: 'Warning', 
            text: "you didn't pay anything", 
            type: 'info',
            confirmButtonColor: "#3da09a"})
    }
    var url = "/transaction/savesalespayment"
    var formData = new FormData()
    formData.append('data', JSON.stringify(change))
    formData.append('LastNumber', salespayment.sequenceNumber())
    var attachment = document.getElementById('uploadFile');
    if (attachment.files.length!=0){
        if (attachment[i].files[0])
        if (!ProActive.checkUploadFilesize(attachment[i].files[0].size, 
            "Attachment file size cannot exceed $size (Row #" (+ i+1) + ")!")) return;
        formData.append("fileUpload", attachment.files[0]);
    }
    // var param = {
    //     Data: change,
    //     LastNumber: salespayment.sequenceNumber()
    // }
    swal({
        title: "Are you sure?",
        text: "You will save this Sales Payment",
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
                    setTimeout(function () {
                        swal({
                            title: "Success!",
                            text: "Data has been saved!",
                            type: "success",
                            confirmButtonColor: "#3da09a"
                        }, function (e) {
                            var x = false
                            if (e) {
                            x = salespayment.printToPdf()
                            }
                            if (x) {
                                location.reload()
                            }
                        });
                    }, 100)
                    model.Processing(false)
                }
            });
            // model.Processing(true)
            // ajaxPost(url, param, function (e) {
            //     setTimeout(function () {
            //         swal({
            //             title: "Success!",
            //             text: "Data has been saved!",
            //             type: "success",
            //             confirmButtonColor: "#3da09a"
            //         }, function (e) {
            //             var x = false
            //             if (e) {
            //                x = salespayment.printToPdf()
            //             }
            //             if (x) {
            //                 location.reload()
            //             }
            //         });
            //     }, 100)
            //     model.Processing(false)
            // })
        } else {
            swal({
                title: "Cancelled", 
                text: "", 
                type: "error",
                confirmButtonColor: "#3da09a"});
        }
    });
 }

salespayment.onChangeCustomerCode = function (value) {
    var data = _.filter(salespayment.dataMasterCustomer(), {
        'CustomerCode': value
    });
    if (data.length == 0) {
        swal({
            title: "Warning!", 
            text: "Choose the Customer Code that have transaction", 
            type: "info",
            confirmButtonColor: "#3da09a"})
        salespayment.resetView()
    } else {
        salespayment.record.CustomerName(data[0].CustomerName)
        salespayment.record.CustomerCode(data[0].CustomerCode)
        salespayment.record.ListDetail([])
        for (var i = 0; i < data.length; i++) {
            salespayment.record.ListDetail.push(salespayment.listDetail({}))
            var dateStr = moment(data[i].DateCreated).format('DD-MMM-YYYY')
            salespayment.record.ListDetail()[i].DatePayment(dateStr)
            salespayment.record.ListDetail()[i].InvNumber(data[i].DocumentNo)
            salespayment.record.ListDetail()[i].Amount(data[i].GrandTotalIDR)
            salespayment.record.ListDetail()[i].AlreadyPaid(data[i].AlreadyPaid)
            salespayment.record.ListDetail()[i].Balance(data[i].GrandTotalIDR - data[i].AlreadyPaid )
        }
        $("[name='pay']").bootstrapSwitch();
        salespayment.maskingMoney()
    }
}

salespayment.maskingMoney = function () {
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

salespayment.selectAll = function () {
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

salespayment.onChangePaid = function (n, value) {
    var index = n.replace("pay", "");
    salespayment.record.ListDetail()[index].Pay(value)
    if (value) {
        $('#recive' + index).removeAttr("disabled");
    } else {
        $('#recive' + index).attr('disabled', 'disabled')
    }
}

salespayment.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
salespayment.createSPFromPOInvSummary = function(){
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    var stat = url.searchParams.get("stat");
    if (num != null){
        // console.log(num)
        if(stat== "PAID"){
            var data = _.filter(salespayment.dataMasterSalesPaymentOriginal(), {
                'CustomerCode': num,
            });
            salespayment.viewList(data[0]._id)
        }else{
            salespayment.createdForm()
            $('.nav-tabs a[href="#CreateSalesPayment"]').tab('show');
            $("#customercode").data("kendoDropDownList").value(num)
            salespayment.onChangeCustomerCode(num)
        }
    }
}
salespayment.processsPlusPlus = function(){
    processs+=1
    // console.log(processs)
    if(processs>=4){
        setTimeout(function(){
            salespayment.createSPFromPOInvSummary()
        },300)
        // setTimeut(purchasepayment.createPPFromPOInvSummary(),10000) 
    }
}
salespayment.init = function () {
    var now = new Date()
    salespayment.getDataPaymentAccount()
    salespayment.getCustomerCode()
    // salespayment.initDate()
    // salespayment.setDate()
    salespayment.getDocumentNumber()
    salespayment.getDataInvoice()
    salespayment.getAllDataSalesPayment(function () {
        salespayment.renderGrid()
    })
    salespayment.getDateNow()
    salespayment.record.DatePosting(moment(new Date()).format('DD-MMM-YYYY'))
}

$(function () {
    $(document).on('change', ':file', function() {
        var input = $(this),
            numFiles = input.get(0).files ? input.get(0).files.length : 1,
            label = input.val().replace(/\\/g, '/').replace(/.*\//, '');
        input.trigger('fileselect', [numFiles, label]);
    });
    $(':file').on('fileselect', function(event, numFiles, label) {

        var input = $(this).parents('.input-group').find(':text'),
            log = numFiles > 1 ? numFiles + ' files selected' : label;

        if( input.length ) {
            input.val(log);
        } else {
            if( log ) alert(log);
        }

    });
    salespayment.init()
    $('#tableSP').on('switchChange.bootstrapSwitch', 'input[name="pay"]', function (event, state) {
        salespayment.onChangePaid(event.currentTarget.id, state)
    });
    $("#textSearch").on("keyup blur change", function () {
        salespayment.filterText();
    });
})

// salespayment.initDate = function () {
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

salespayment.onChangeDateStart = function(val){
    if (val.getTime()>salespayment.DateEnd().getTime()){
        salespayment.DateEnd(val)
    }
}

salespayment.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["StoreLocationName", "DocumentNumber", "CustomerName"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridListSalesPayment").data("kendoGrid").dataSource.query({ filter: filter });
}