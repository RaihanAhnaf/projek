model.Processing(false)
var firstamount = {}

firstamount.docnofromdb = ko.observable()
firstamount.showListJournal = ko.observable(false)
firstamount.showListDraft = ko.observable(false)
firstamount.PrintiIsActive = ko.observable(false)
firstamount.dataJurnal = ko.observableArray([])
firstamount.cashNameJournal = ko.observable("")
firstamount.balanceData = ko.observable(false)
firstamount.documentNoStr = ko.observable("")
firstamount.dataMasterAccount = ko.observableArray([])
firstamount.accountNumberForDropDown = ko.observableArray([])
firstamount.documentNumberLastRow = ko.observable("")
firstamount.textListDraftDetail = ko.observable("")
firstamount.typeFilter = ko.observable("")
firstamount.titleListJournal = ko.observable("")
firstamount.dataMasterJournal = ko.observableArray([])
firstamount.dataDraftJournal = ko.observableArray([])
firstamount.dateStartFilter = ko.observable()
firstamount.dateEndFilter = ko.observable("")
firstamount.showIconAttachment = ko.observable(false)
firstamount.showPdf = ko.observable(false)
firstamount.typeFilterForJournal = ko.observable('')
firstamount.idForPosting = ko.observable()
firstamount.createDate = ko.observable("")
firstamount.refreshDraftwithFilter = ko.observable(false)
firstamount.lastDate = ko.observable()
firstamount.DatePageBar = ko.observable()

firstamount.typeAllJournal = [{
    "value": "CashIn",
    "text": "Cash In"
}, {
    "value": "CashOut",
    "text": "Cash Out"
}, {
    "value": "General",
    "text": "General"
}]

firstamount.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    firstamount.DatePageBar(page)
}

firstamount.cashNameJournal.subscribe(function (e) {
    $('#textSearch').val("")
})

firstamount.newRecord = function () {
    var page = {
        ID: "",
        IdJournal: "",
        PostingDate: "",
        CreateDate: "",
        DateStr: "",
        User: "",
        Journal_Type: "",
        ListDetail: [],
        Status: "",
    }
    page.ListDetail.push(firstamount.listDetailJournal())
    return page
}

firstamount.listDetailJournal = function () {
    return {
        Id: ko.observable(''),
        No: ko.observable(0),
        Journal_Type: ko.observable(''),
        PostingDate: ko.observable(''),
        DateStr: ko.observable(''),
        DocumentNumber: ko.observable(''),
        Acc_Code: ko.observable(''),
        Acc_Name: ko.observable(''),
        Debet: ko.observable(0),
        Credit: ko.observable(0),
        CreditDisable: ko.observable(false),
        Description: ko.observable(''),
        Attachment: ko.observable(''),
        User: ko.observable(''),
    }
}

firstamount.record = ko.mapping.fromJS(firstamount.newRecord())

firstamount.getData = function () {
    model.Processing(true)

    ajaxPost('/transaction/getaccount', {}, function (res) {

        if (res.Total === 0) {
            swal({
                title:"Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        firstamount.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        firstamount.accountNumberForDropDown(DataAccount)

        model.Processing(false)
    })
}

firstamount.totalDebet = ko.computed(function () {
    if (!_.isArray(firstamount.record.ListDetail())) return 0;
    var totalDebet = _.sumBy(firstamount.record.ListDetail(), function (v) {
        var debet = kendo.toString(v.Debet(), "n");
        var debetInt = Number(debet.replace(/[^0-9\.]+/g, ""));
        return debetInt
    })
    return kendo.toString(totalDebet, "n");
})

firstamount.totalCredit = ko.computed(function () {
    if (!_.isArray(firstamount.record.ListDetail())) return 0;
    var totalCredit = _.sumBy(firstamount.record.ListDetail(), function (v) {
        var credit = kendo.toString(v.Credit(), "n");
        var creditInt = Number(credit.replace(/[^0-9\.]+/g, ""));
        return creditInt
    })
    return kendo.toString(totalCredit, "n");
})

firstamount.balance = ko.computed(function () {
    var totalDebet = Number(firstamount.totalDebet().replace(/[^0-9\.]+/g, ""));
    var totalCredit = Number(firstamount.totalCredit().replace(/[^0-9\.]+/g, ""));
    var total = totalDebet - totalCredit

    if (total >= 0) {
        var TotString = kendo.toString(total, "n");
        return TotString;
    } else {
        var TotminString = kendo.toString(Math.abs(total), "n");
        return "(" + TotminString + ")";
    }

})

firstamount.onChangeAccountNumber = function (value, index) {
    findaccount = _.find(firstamount.dataMasterAccount(), {
        ACC_Code: value
    })
    firstamount.record.ListDetail()[index].Acc_Name(findaccount.Account_Name);
}

firstamount.addContent = function () {
    var documentcode = ""
    var journaltype = ""
    firstamount.record.ListDetail.push(firstamount.listDetailJournal())
    firstamount.record.User(userinfo.usernameh());
    firstamount.record.CreateDate(new Date());
    firstamount.record.PostingDate(new Date());
    firstamount.record.DateStr(moment().format("DD MMMM YYYY"));
    firstamount.record.Status("posting");
    firstamount.record.ListDetail()[firstamount.record.ListDetail().length - 1].No(firstamount.record.ListDetail().length);
    firstamount.record.ListDetail()[firstamount.record.ListDetail().length - 1].PostingDate(new Date());
    firstamount.record.ListDetail()[firstamount.record.ListDetail().length - 1].User(userinfo.usernameh());

    switch (firstamount.cashNameJournal()) {
    case "Cash In":
        documentcode = "BBM"
        journaltype = "CashIn"
        break
    case "Cash Out":
        documentcode = "BKK"
        journaltype = "CashOut"
        break
    default:
        documentcode = "GEM"
        journaltype = "General"
    }
    var number = 0
    if (firstamount.record.ListDetail().length < 2) {

        number = 1

    } else {

        if (firstamount.balance() == "0.00" && firstamount.totalDebet() != "0.00") {
            var jr = ko.mapping.toJS(firstamount.record)
            var i = firstamount.record.ListDetail().length - 2
            var j = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 4)
            parseInt(j)
            var no = 0
            if (j < 9) {
                no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 1)
            }
            if (j >= 10 && j <= 100) {
                no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 2)
            }
            if (j >= 100 && j <= 1000) {
                no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 3)
            }
            if (j >= 1000 && j <= 9999) {
                no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 4)
            }
            var Nom = parseInt(no)
            number = Nom + 1
        } else if (firstamount.totalDebet() == "0.00" || firstamount.totalCredit() == "0.00") {
            number = firstamount.documentNumberLastRow()
        } else {
            number = firstamount.documentNumberLastRow()
        }
    }
    var beforeNumber = "000"
    if (number >= 10 && number <= 100) {
        beforeNumber = "00"
    }
    if (number >= 100 && number <= 1000) {
        beforeNumber = "0"
    }
    if (number >= 1000 && number <= 9999) {
        beforeNumber = ""
    }
    firstamount.record.Journal_Type(journaltype);
    firstamount.record.ListDetail()[firstamount.record.ListDetail().length - 1].Id(firstamount.makeId());
    firstamount.record.ListDetail()[firstamount.record.ListDetail().length - 1].DocumentNumber(documentcode + "/" + moment().format("DDMMYY") + "/temp-" + userinfo.usernameh() + "/" + beforeNumber + number);
    firstamount.record.ListDetail()[firstamount.record.ListDetail().length - 1].Journal_Type(journaltype);
    firstamount.maskingMoney()
    firstamount.documentNumberLastRow(number)
}

firstamount.getLastDocNumber = function () {
    switch (firstamount.cashNameJournal()) {
    case "Journal Cash In":
        journaltype = "CashIn"
        break
    case "Journal Cash Out":
        journaltype = "CashOut"
        break
    default:
        journaltype = "General"
    }
    var param = {
        Type: journaltype
    }
    ajaxPost('/transaction/getdocumentnumber', param, function (res) {
        if (res.Data.length != 0) {
            var docno = res.Data[0].DocNo
            var no = docno.substr(docno.length - 1)
            firstamount.docnofromdb(parseInt(no))
        }

    })
}

firstamount.addNewRow = function () {
    if (firstamount.cashNameJournal() == "") {
        swal({
            title:"Error!",
            text: "Please Choose Journal Type for the first",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
        return
    }
    if (firstamount.record.ListDetail().length > 0) {
        var jr = ko.mapping.toJS(firstamount.record)
        var i = firstamount.record.ListDetail().length - 1
        var no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 1);
        var index = firstamount.record.ListDetail().length;
        firstamount.addContent()
        return
    }
    firstamount.addContent()
}

firstamount.maskingMoney = function () {
    $('.currency').inputmask("numeric", {
        radixPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        rightAlign: false,
    });
}

firstamount.removeRow = function () {
    firstamount.record.ListDetail.remove(this)
    if (firstamount.record.ListDetail().length != 0) {
        var jr = ko.mapping.toJS(firstamount.record)
        var i = firstamount.record.ListDetail().length - 1
        var j = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 4)
        parseInt(j)
        var no = 0
        if (j < 9) {
            no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 1)
        }
        if (j >= 10 && j <= 100) {
            no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 2)

        }
        if (j >= 100 && j <= 1000) {
            no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 3)
        }
        if (j >= 1000 && j <= 9999) {
            no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 4)
        }
        var Nom = parseInt(no)
        firstamount.documentNumberLastRow(Nom)
    }
}

firstamount.changetoCashIn = function () {
    if (firstamount.record.ListDetail() != 0) {
        swal({
            title: "Are you sure?",
            text: "You will remove all your data now!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3da09a",
            confirmButtonText: "Yes!",
            cancelButtonText: "No!",
            closeOnConfirm: true,
            closeOnCancel: true
        }, function () {
            firstamount.record.ListDetail([]);
            firstamount.cashNameJournal("Cash In")
        });
    } else {
        firstamount.record.ListDetail([]);
        firstamount.cashNameJournal("Cash In")
    }
}

firstamount.changetoCashOut = function () {
    if (firstamount.record.ListDetail() != 0) {
        swal({
            title: "Are you sure?",
            text: "You will remove all your data now!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3da09a",
            confirmButtonText: "Yes!",
            cancelButtonText: "No!",
            closeOnConfirm: true,
            closeOnCancel: true
        }, function () {
            firstamount.record.ListDetail([]);
            firstamount.cashNameJournal("Cash Out")
        });
    } else {
        firstamount.record.ListDetail([]);
        firstamount.cashNameJournal("Cash Out")
    }

}

firstamount.changetoGeneral = function () {
    if (firstamount.record.ListDetail() != 0) {
        swal({
            title: "Are you sure?",
            text: "You will remove all your data now!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3da09a",
            confirmButtonText: "Yes!",
            cancelButtonText: "No!",
            closeOnConfirm: true,
            closeOnCancel: true
        }, function () {
            firstamount.record.ListDetail([]);
            firstamount.cashNameJournal("General Journal")
        });
    } else {
        firstamount.record.ListDetail([]);
        firstamount.cashNameJournal("General Journal")
    }

}

firstamount.makeId = function () {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 10; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}
firstamount.resetRow = function () {
    firstamount.record.ListDetail([]);
}

firstamount.saveData = function () {
    var payload = ko.mapping.toJS(firstamount.record)
    var formData = new FormData()
    if (payload.ListDetail.length == 0) {
        return swal({
            title:"Warning!",
            text: "None data for save",
            type: "info",
            confirmButtonColor:"#3da09a"})
    }
    var i = firstamount.record.ListDetail().length - 1
    var accname = payload.ListDetail[i].Acc_Code
    if (accname == "") {
        return swal({
            title:"Warning!", 
            text:"Please check Account Code",
            type: "info",
            confirmButtonColor:"#3da09a"})
    } else if (firstamount.totalDebet() <= 0 || firstamount.totalCredit() <= 0 || firstamount.balance() != 0) {
        return swal({
            title:"Warning!",
            text: "Please Balance the Journal for the first", 
            type:"info",
            confirmButtonColor:"#3da09a"})
    }


    _.each(payload.ListDetail, function (v, i) {
        if (v.Acc_Name == "") {}
        payload.ListDetail[i].Acc_Code = parseInt(v.Acc_Code)

        if (v.Credit == 0) {
            payload.ListDetail[i].Credit = parseFloat(v.Credit)
        }
        if (v.Debet == 0) {
            payload.ListDetail[i].Debet = parseFloat(v.Debet)
        }
        if (v.Credit != 0) {
            payload.ListDetail[i].Credit = Number(v.Credit.replace(/[^0-9\.]+/g, ""));
        }
        if (v.Debet != 0) {
            payload.ListDetail[i].Debet = Number(v.Debet.replace(/[^0-9\.]+/g, ""));
        }
        payload.ListDetail[i].DateStr = moment(payload.ListDetail[i].PostingDate).format("DD MMMM YYYY")
        payload.PostingDate = payload.ListDetail[i].PostingDate
        payload.DateStr = moment(payload.PostingDate).format("DD MMMM YYYY")
    })
    formData.append('data', JSON.stringify(payload))
    var i
    var attachment = document.getElementsByClassName('upload');
    for (i = 0; i < attachment.length; i++) {

        formData.append("fileUpload" + i, attachment[i].files[0]);
    }

    var url = "/transaction/savejournalfirstamount"
    swal({
        title: "Are you sure?",
        text: "You will Posting this journal",
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
                    firstamount.getLastDocNumber()
                    firstamount.resetRow()
                }
            });
            model.Processing(false)
        } else {
            swal({
                title:"Cancelled",
                type:"error",
                confirmButtonColor:"#3da09a"});
        }
    });

}

firstamount.init = function () {
    firstamount.getData();
    firstamount.cashNameJournal("Cash In")
    firstamount.getDateNow()
}

$(function () {
    firstamount.init()
    firstamount.record.ListDetail([]);
    firstamount.getLastDocNumber()
    $(".btn").mouseup(function () {
        $(this).blur();
    })
})