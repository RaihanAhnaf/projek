model.Processing(false)

var journal = {}
journal.docnofromdb = ko.observable()
journal.showListJournal = ko.observable(false)
journal.showListDraft = ko.observable(false)
journal.PrintiIsActive = ko.observable(false)
journal.dataJurnal = ko.observableArray([])
journal.cashNameJournal = ko.observable("")
journal.balanceData = ko.observable(false)
journal.documentNoStr = ko.observable("")
journal.dataMasterAccount = ko.observableArray([])
journal.accountNumberForDropDown = ko.observableArray([])
journal.documentNumberLastRow = ko.observable("")
journal.textListDraftDetail = ko.observable("")
journal.typeFilter = ko.observable("")
journal.titleListJournal = ko.observable("")
journal.dataMasterJournal = ko.observableArray([])
journal.dataDraftJournal = ko.observableArray([])
journal.dateStartFilter = ko.observable()
journal.dateEndFilter = ko.observable("")
journal.showIconAttachment = ko.observable(false)
journal.showPdf = ko.observable(false)
journal.typeFilterForJournal = ko.observable('')
journal.idForPosting = ko.observable()
journal.createDate = ko.observable("")
journal.refreshDraftwithFilter = ko.observable(false)
journal.lastDate = ko.observable()
journal.TabsIndexPosition = ko.observable()
journal.DatePageBar = ko.observable()
journal.textSearch = ko.observable()
journal.ShowFilterJournal = ko.observable(true)
journal.valueDepartment = ko.observable("")
journal.valueSales = ko.observable("")
journal.nameSales = ko.observable("")
journal.idJournalFromCoa = ko.observable("")
journal.idJournalFromCoaForUnapply = ko.observable("")
journal.showUnApply = ko.observable(false)
journal.statusJournal = ko.observable("")
journal.dateStart = ko.observable(new Date)
journal.dateEnd = ko.observable(new Date)
journal.typeAllJournal = [{
    "value": "CashIn",
    "text": "Cash In"
}, {
    "value": "CashOut",
    "text": "Cash Out"
}, {
    "value": "General",
    "text": "General"
}]

journal.cashNameJournal.subscribe(function (e) {
    $('#textSearch').val("")
})

journal.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    journal.DatePageBar(page)
}

journal.newRecord = function () {
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
        Department: "",
        SalesCode: "",
        SalesName: ""
    }
    page.ListDetail.push(journal.listDetailJournal())
    return page
}

journal.listDetailJournal = function () {
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
        Department: ko.observable(''),
        SalesCode: ko.observable(''),
        SalesName: ko.observable('')
    }
}

journal.record = ko.mapping.fromJS(journal.newRecord())

journal.getLastDocNumber = function () {
    switch (journal.cashNameJournal()) {
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
            journal.docnofromdb(parseInt(no))
        }

    })
}
journal.getData = function () {
    model.Processing(true)

    ajaxPost('/transaction/getaccount', {}, function (res) {

        if (res.Total === 0) {
            swal("Error!", res.Message, "error")
            return
        }
        journal.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        journal.accountNumberForDropDown(DataAccount)

        model.Processing(false)
    })
}

journal.getDataJournal = function (callback) {
    var s = $('#dateStartAllJournal').data('kendoDatePicker').value();
    var e = $('#dateEndAllJournal').data('kendoDatePicker').value();
    var text = ""
    if ($("#filterDescription").val() != "") {
        text = $("#filterDescription").val()
    } else {
        text = journal.textSearch()
    }
    var param = {
        Type: journal.typeFilterForJournal(),
        TextSearch: text,
        DateStart: moment(s).format('YYYY-MM-DD'),
        DateEnd: moment(e).format('YYYY-MM-DD'),
        Filter: journal.refreshDraftwithFilter(),
        IdJournal: journal.idJournalFromCoa()
    }
    model.Processing(true)
    ajaxPost('/transaction/getdatajournal', param, function (res) {
        model.Processing(false)
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        journal.dataMasterJournal(res.Data.Data1)
        if (res.Data.Data2.length != 0) {
            journal.statusJournal(res.Data.Data2[0].Status)
        }
        for (var i = 0; i < journal.dataMasterJournal().length; i++) {
            journal.dataMasterJournal()[i].DownloadLink = moment(journal.dataMasterJournal()[i].PostingDate).format('YYYYMM') + '/' + journal.dataMasterJournal()[i].Attachment
        }
        journal.idJournalFromCoa("")

        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

journal.downloadAttachment = function (e, f) {
    //console.log(e, f)
    if (f != "") {
        var pom = document.createElement('a');
        pom.setAttribute('href', "/res/docs/" + e);
        pom.setAttribute('download', e);
        pom.click();

    }
    if (f == "") {
        return swal({
            title: 'Warning!',
            text: 'this data has no attachment',
            type: 'info',
            confirmButtonColor: '#3da09a',
        })
    }
    if (f == "BEGIN" || f == "INVOICE") {
        return swal({
            title: 'Warning!',
            text: 'this is begining data',
            type: 'info',
            confirmButtonColor: '#3da09a',
        })
    }
}

journal.renderGrid = function () {
    var data = journal.dataMasterJournal();
    if (typeof $('#listJournal').data('kendoGrid') !== 'undefined') {
        $('#listJournal').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 25,
            group: [{
                field: "Journal_Type"
            }, {
                field: "DateStr"
            }],
            sort: ({
                field: "DocumentNumber",
                dir: "desc"
            }),
        }))
        return
    }

    var columns = [{
        field: 'DocumentNumber',
        title: 'DocumentNumber',
        width: 100
    }, {
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 70
    }, {
        field: 'Acc_Name',
        title: 'Acc Name',
        width: 100
    }, {
        field: 'Department',
        title: 'Department',
        width: 80
    }, {
        title: 'Sales',
        width: 80,
        template: function (d) {
            if (d.SalesName == "" || d.SalesName == undefined) return "--"
            else return d.SalesName
        }
    }, {
        field: 'Debet',
        title: 'Debit',
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        template: "#=ChangeToRupiah(Debet)#"
    }, {
        field: 'Credit',
        title: 'Credit',
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        template: "#=ChangeToRupiah(Credit)#"
    }, {
        field: 'Description',
        title: 'Description',
        width: 100
    }, {
        field: 'User',
        title: 'User',
        width: 85
    }, {
        title: 'Attachment',
        width: 70,
        attributes: {
            "class": "centerAction",
        },
        template: function (e) {
            var kosong = ""
            if (e.Attachment != "" && e.Attachment != "BEGIN") {
                return "<a href=\"javascript:journal.downloadAttachment('" + e.DownloadLink + "', '" + e.Attachment + "')\" data-target=\".downloadAttachment\" data-bind=\"attr: {id: 'downloadlink' + $index()}\" data-backdrop=\"static\" class=\"btn btn-xs btn-primary\"><i class=\"glyphicon glyphicon-paperclip\"></i></a>&nbsp;"
            }
            return "<a href=\"javascript:journal.downloadAttachment('" + kosong + "', '" + kosong + "')\" data-target=\".downloadAttachment\" data-bind=\"attr: {id: 'downloadlink' + $index()}\" data-backdrop=\"static\" class=\"btn btn-xs btn-danger\"><i class=\"glyphicon glyphicon-paperclip\"></i></a>&nbsp;"
        }
    }, {
        field: "Journal_Type",
        template: "",
        hidden: true,
        groupHeaderTemplate: "Journal Type : #= value #"
    }, {
        field: "DateStr",
        template: "",
        hidden: true,
        groupHeaderTemplate: "Date : #= value #"
    },]

    $('#listJournal').kendoGrid({
        dataSource: {
            data: data,
            pageSize: 25,
            group: [{
                field: "Journal_Type"
            }, {
                field: "DateStr"
            }],
            sort: ({
                field: "DocumentNumber",
                dir: "desc"
            }),
        },
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        height: 500,
        sortable: true,
        scrollable: true,
        columns: columns,
    })

}
journal.addContent = function () {
    var documentcode = ""
    var journaltype = ""
    journal.record.ListDetail.push(journal.listDetailJournal())
    journal.record.User(userinfo.usernameh());
    journal.record.CreateDate(new Date());
    journal.record.PostingDate(new Date());
    journal.record.DateStr(moment().format("DD MMMM YYYY"));
    journal.record.Department(journal.valueDepartment())
    journal.record.Status("draft");
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].No(journal.record.ListDetail().length);
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].PostingDate(new Date());
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].User(userinfo.usernameh());

    switch (journal.cashNameJournal()) {
        case "Journal Cash In":
            documentcode = "BBM"
            journaltype = "CashIn"
            break
        case "Journal Cash Out":
            documentcode = "BKK"
            journaltype = "CashOut"
            break
        default:
            documentcode = "GEM"
            journaltype = "General"
    }
    var number = 0
    if (journal.record.ListDetail().length < 2) {
        number = 1
    } else {
        if (journal.balance() == "0.00" && journal.totalDebet() != "0.00") {
            var jr = ko.mapping.toJS(journal.record)
            var i = journal.record.ListDetail().length - 2
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
        } else if (journal.totalDebet() == "0.00" || journal.totalCredit() == "0.00") {
            number = journal.documentNumberLastRow()
        } else {
            number = journal.documentNumberLastRow()
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

    var datecodenumber = ""

    if (userinfo.rolenameh() == "supervisor") {
        datecodenumber = moment(journal.record.ListDetail()[journal.record.ListDetail().length - 1].PostingDate()).format("DDMMYY")
    } else {
        datecodenumber = moment().format("DDMMYY")
    }

    journal.record.Journal_Type(journaltype);
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].Id(journal.makeId());
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].DocumentNumber(documentcode + "/" + datecodenumber + "/temp-" + userinfo.usernameh() + "/" + beforeNumber + number);
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].Journal_Type(journaltype);
    journal.maskingMoney()
    journal.documentNumberLastRow(number)
    journal.getLastDate()
    if (userinfo.rolenameh() != "supervisor" && userinfo.rolenameh() != "administrator") {
        journal.disableDate()
    }
}
journal.getLastDate = function () {
    var url = "/closing/getlastdate"
    ajaxPost(url, {}, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        if (res.Data.length == 0) {
            journal.lastDate(new Date(01, 0, 01))
        } else {
            var years = moment(res.Data[0].lastdate).year()
            var month = moment(res.Data[0].lastdate).month() + 1
            journal.lastDate(new Date(years, month, 01))
        }

    })
}
journal.maskingMoney = function () {
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
journal.disableDate = function () {
    $(".datePosting").each(function () {
        $(this).find(':input').attr('disabled', 'disabled').addClass("k-state-disabled")
        $(this).find('.k-select').remove()
        $(this).find('.datepicker').unwrap()
    })
}
journal.addNewRow = function () {
    if (journal.cashNameJournal() == "Journal") {
        swal({
            title: "Error!",
            text: "Please Choose Journal Type for the first",
            type: "error",
            confirmButtonText: "OK"
        })
    } else if (journal.record.ListDetail().length > 0) {
        var jr = ko.mapping.toJS(journal.record)
        var i = journal.record.ListDetail().length - 1
        var no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 1);
        var index = journal.record.ListDetail().length;
        if (jr.ListDetail[i].Debet == 0 && jr.ListDetail[i].Credit == 0) {
            swal({
                title: "Error!",
                text: "Please check debit or credit journal",
                type: "error",
                confirmButtonColor: '#3da09a',
                confirmButtonText: "OK"
            })
        } else if (jr.ListDetail[i].Acc_Code == "") {
            swal({
                title: "Error!",
                text: "Please check Account Code",
                type: "error",
                confirmButtonColor: '#3da09a',
                confirmButtonText: "OK"
            })
        } else if (journal.record.ListDetail().length % 2 != 0) {
            journal.addContent()
            setTimeout(function () {
                $('#debit_' + index).prop("disabled", true);
            }, 100);

        } else {
            if (journal.balance() != 0) {
                swal({
                    title: "Error!",
                    text: "Please Balance the Journal for the first",
                    type: "error",
                    confirmButtonColor: '#3da09a',
                    confirmButtonText: "OK"
                })
                return
            } else {
                (journal.totalDebet() == 0 && journal.totalCredit() == 0 || journal.balance() != 0)
                journal.addContent()
                setTimeout(function () {
                    $('#credit_' + index).prop("disabled", true);
                }, 100);
            }
        }
    } else {
        journal.addContent()
        $('#credit_0').prop("disabled", true);
    }
}

journal.onChangeAccountNumber = function (value, index) {
    findaccount = _.find(journal.dataMasterAccount(), {
        ACC_Code: value
    })
    journal.record.ListDetail()[index].Acc_Name(findaccount.Account_Name);
}

journal.onChangeDate = function (value, index) {
    switch (journal.cashNameJournal()) {
        case "Journal Cash In":
            documentcode = "BBM"
            journaltype = "CashIn"
            break
        case "Journal Cash Out":
            documentcode = "BKK"
            journaltype = "CashOut"
            break
        default:
            documentcode = "GEM"
            journaltype = "General"
    }
    var number = 0
    if (journal.record.ListDetail().length < 2) {
        number = 1
    } else {
        if (journal.balance() == "0.00" && journal.totalDebet() != "0.00") {
            var jr = ko.mapping.toJS(journal.record)
            var i = journal.record.ListDetail().length - 2
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
        } else if (journal.totalDebet() == "0.00" || journal.totalCredit() == "0.00") {
            number = journal.documentNumberLastRow()
        } else {
            number = journal.documentNumberLastRow()
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

    var datecodenumber = ""

    if (userinfo.rolenameh() == "supervisor") {
        datecodenumber = moment(value).format("DDMMYY")
    } else {
        datecodenumber = moment().format("DDMMYY")
    }
    journal.record.ListDetail()[journal.record.ListDetail().length - 1].DocumentNumber(documentcode + "/" + datecodenumber + "/temp-" + userinfo.usernameh() + "/" + beforeNumber + number);
}

journal.removeRow = function () {
    journal.record.ListDetail.remove(this)
    if (journal.record.ListDetail().length != 0) {
        var jr = ko.mapping.toJS(journal.record)
        var i = journal.record.ListDetail().length - 1
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
        journal.documentNumberLastRow(Nom)
    }
}

journal.changetoCashIn = function () {
    $("#filtered").prop('disabled', true);
    journal.record.ListDetail([])
    journal.cashNameJournal("Journal Cash In")
    journal.addNewRow()
    journal.ShowFilterJournal(false)
    journal.departmentDropdown()
    // var dropdownlist = $("#departmenDropdown").data("kendoDropDownList");
    // dropdownlist.value("");
    // dropdownlist.trigger("change"); 

}

journal.changetoCashOut = function () {
    $("#filtered").prop('disabled', true);
    journal.record.ListDetail([]);
    journal.cashNameJournal("Journal Cash Out")
    journal.addNewRow()
    journal.ShowFilterJournal(false)
    journal.departmentDropdown()
    journal.salesDropdown()
    // var dropdownlist = $("#departmenDropdown").data("kendoDropDownList");
    // dropdownlist.value("");
    // dropdownlist.trigger("change"); 
}

journal.changetoGeneral = function () {
    $("#filtered").prop('disabled', true);
    journal.record.ListDetail([]);
    journal.cashNameJournal("General Journal")
    journal.addNewRow()
    journal.ShowFilterJournal(false)
    journal.departmentDropdown()
    journal.salesDropdown()
    // var dropdownlist = $("#departmenDropdown").data("kendoDropDownList");
    // dropdownlist.value("");
    // dropdownlist.trigger("change"); 

}

journal.totalDebet = ko.computed(function () {
    if (!_.isArray(journal.record.ListDetail())) return 0;
    var totalDebet = _.sumBy(journal.record.ListDetail(), function (v) {
        var debet = kendo.toString(v.Debet(), "n");
        var debetInt = Number(debet.replace(/[^0-9\.]+/g, ""));
        return debetInt
    })
    return kendo.toString(totalDebet, "n");
})

journal.totalCredit = ko.computed(function () {
    if (!_.isArray(journal.record.ListDetail())) return 0;
    var totalCredit = _.sumBy(journal.record.ListDetail(), function (v) {
        var credit = kendo.toString(v.Credit(), "n");
        var creditInt = Number(credit.replace(/[^0-9\.]+/g, ""));
        return creditInt
    })
    return kendo.toString(totalCredit, "n");
})

journal.balance = ko.computed(function () {
    var totalDebet = Number(journal.totalDebet().replace(/[^0-9\.]+/g, ""));
    var totalCredit = Number(journal.totalCredit().replace(/[^0-9\.]+/g, ""));
    var total = totalDebet - totalCredit

    if (total >= 0) {
        var TotString = kendo.toString(total, "n");
        return TotString;
    } else {
        var TotminString = kendo.toString(Math.abs(total), "n");
        return "(" + TotminString + ")";
    }
})

journal.debitChange = function (d, i) {
    if (parseFloat(d()) > 0) {
        $("#credit_" + i).prop('enable', true);
    } else {
        $("#credit_" + i).prop('disabled', false);
    }


}

journal.creditChange = function (d, i) {
    if (parseFloat(d()) > 0) {
        $("#debit_" + i).prop('disabled', true);
    } else {
        $("#debit_" + i).prop('enable', false);
    }

}

journal.saveData = function () {
    var payload = ko.mapping.toJS(journal.record)
    var formData = new FormData()
    if (payload.ListDetail.length == 0) {
        return swal({
            title: "Warning!",
            text: "None data for save",
            type: "info",
            confirmButtonColor: '#3da09a'
        })
    }
    if (journal.valueDepartment() == "") {
        return swal({
            title: "Warning!",
            text: "Department not selected",
            type: "info",
            confirmButtonColor: '#3da09a'
        })
    }

    if (journal.valueSales() == "" && (journal.TabsIndexPosition() == 1 || journal.TabsIndexPosition() == 2)) {
        return swal({
            title: "Warning!",
            text: "Sales not selected",
            type: "info",
            confirmButtonColor: '#3da09a'
        })
    }

    payload.Department = journal.valueDepartment()
    payload.SalesCode = journal.valueSales()
    payload.SalesName = journal.nameSales()
    var i = journal.record.ListDetail().length - 1
    var accname = payload.ListDetail[i].Acc_Code
    if (accname == "") {
        return swal({
            title: "Warning!",
            text: "Please check Account Code",
            type: "info",
            confirmButtonColor: '#3da09a'
        })
    } else if (journal.totalDebet() <= 0 || journal.totalCredit() <= 0 || journal.balance() != 0) {
        return swal({
            title: "Warning!",
            text: "Please Balance the Journal for the first",
            type: "info",
            confirmButtonColor: '#3da09a'
        })
    }

    _.each(payload.ListDetail, function (v, i) {
        if (v.Acc_Name == "") { }
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
        payload.ListDetail[i].Department = journal.valueDepartment()
        payload.ListDetail[i].SalesCode = journal.valueSales()
        payload.ListDetail[i].SalesName = journal.nameSales()
        payload.PostingDate = payload.ListDetail[i].PostingDate
        payload.DateStr = moment(payload.PostingDate).format("DD MMMM YYYY")
    })
    formData.append('data', JSON.stringify(payload))
    var i;
    var attachment = document.getElementsByClassName('upload');
    for (i = 0; i < attachment.length; i++) {
        if (attachment[i].files[0])
            if (!ProActive.checkUploadFilesize(attachment[i].files[0].size,
                "Attachment file size cannot exceed $size (Row #" + (i + 1) + ")!")) return;
        formData.append("fileUpload" + i, attachment[i].files[0]);
    }

    var url = "/transaction/savejournal"
    swal({
        title: "Are you sure?",
        text: "You will submit this journal",
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
                                confirmButtonColor: "#3da09a",
                            }, function () {
                                // location.reload()
                                $('.nav-tabs a[href="#listofdraft"]').tab('show')
                                journal.getLastDocNumber()
                                journal.getDataDraft(function () {
                                    journal.renderGridDraft()
                                })
                                journal.renderGridDraftDetail()
                                // var dropdownlist = $("#departmenDropdown").data("kendoDropDownList");
                                // dropdownlist.value("");
                                // dropdownlist.trigger("change"); 
                            });

                        }, 100)

                    }
                }
            });
            model.Processing(false)
            if (journal.PrintiIsActive() == true) {
                setTimeout(journal.printJournal(), 1000)
                //console.log("go to function print")
            } else {
                journal.comeToListofDraft();
            }
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });

}
journal.makeId = function () {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 10; i++)
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}
journal.resetRow = function () {
    journal.record.ListDetail([]);
}
journal.postAndPrint = function () {
    // journal.saveData()
    // journal.PrintiIsActive(true)
}
journal.backToJournal = function () {
    $("#filtered").prop('disabled', true);
    journal.showListJournal(false)
    journal.record.ListDetail([]);
}
journal.backToJournal2 = function () {
    $("#filtered").prop('disabled', true);
    journal.refreshDraftwithFilter(false)
    journal.showListDraft(false)
    journal.record.ListDetail([]);
}
journal.comeToListJournal = function () {
    $("#filtered").prop('disabled', false);
    journal.showListJournal(true)
    journal.refresh()
    journal.ShowFilterJournal(true)
    $('#textSearch').val("")
}

journal.tbodyForPrint = function () {
    $('#JournalTable2 tbody').html('')
    var payload = ko.mapping.toJS(journal.record)
    $.each(payload.ListDetail, function (i, d) {
        var attachment = document.getElementsByClassName('upload');
        var row = '<tr>'
        row += '<td>' + d.No + '</td>'
        row += '<td>' + moment(d.PostingDate).format("DD MMM YYYY HH:mm") + '</td>'
        row += '<td>' + d.DocumentNumber + '</td>'
        row += '<td>' + d.Acc_Code + '</td>'
        row += '<td>' + d.Acc_Name + '</td>'
        row += '<td>' + d.Debet + '</td>'
        row += '<td>' + d.Credit + '</td>'
        row += '<td>' + d.Description + '</td>'
        if (attachment[i].files.length > 0) {
            //console.log(attachment)
            row += '<td>' + attachment[i].files.length + "files" + "(" + attachment[i].files[0].name + ")" + '</td>'
        } else {
            row += '<td>' + "nothing file uploaded" + '</td>'
        }
        row += '</tr>';
        $('#JournalTable2 tbody').append(row);
    })
}
journal.backFromPrintDraft = function () {
    journal.showPdf(false)
}
journal.printAllJournal = function () {
    var url = "/transaction/printjournaltopdf"
    var s = $('#dateStartAllJournal').data('kendoDatePicker').value();
    var e = $('#dateEndAllJournal').data('kendoDatePicker').value();
    var param = {
        Type: journal.typeFilterForJournal(),
        DateStart: moment(s).format('YYYY-MM-DD'),
        DateEnd: moment(e).format('YYYY-MM-DD'),
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/journal/' + res.Data, '_blank');
        tabOrWindow.focus();
    })
}
journal.printJournal = function () {
    if (journal.record.ListDetail().length == 0) {
        return swal({
            title: "Warning!",
            text: "Journal is empty",
            type: "info",
            confirmButtonColor: "#3da09a",
        })
    }
    var i;
    var attachment = document.getElementsByClassName('upload');
    for (i = 0; i < attachment.length; i++) {
        if (attachment[i].files.length == 0) {
            return swal({
                title: "Warning!",
                text: "Please upload file",
                type: "info",
                confirmButtonColor: "#3da09a",
            })
        }
    }
    var payload = ko.mapping.toJS(journal.record)
    $.each(payload.ListDetail, function (i, d) {
        payload.ListDetail[i].Acc_Code = parseInt(d.Acc_Code)
        payload.ListDetail[i].Debet = parseFloat(d.Debet)
        payload.ListDetail[i].Credit = parseFloat(d.Credit)
        payload.ListDetail[i].Attachment = attachment[i].files.length + " files" + "(" + attachment[i].files[0].name + ")"
    })
    var param = {
        DateStr: payload.DateStr,
        Journal_Type: payload.Journal_Type,
        Status: payload.Status,
        User: payload.User,
        ListDetail: payload.ListDetail,
        TotalDebet: journal.totalDebet(),
        TotalCredit: journal.totalCredit(),
        Balance: journal.balance()
    }
    model.Processing(true)
    var url = "/transaction/printdraftjournaltopdf"
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/listjournal/' + res, '_blank');
        tabOrWindow.focus();
    })
}

journal.printJournalDraft = function () {
    var dataGrid = $('#listDraftDetail').data('kendoGrid').dataSource.options.data
    if (dataGrid.length == 0) {
        return swal({
            title: "Warning!",
            text: "Journal is empty",
            type: "info",
            confirmButtonColor: "#3da09a",
        })
    }
    var param = {
        ID: journal.idForPosting(),
    }
    model.Processing(true)
    var url = "/transaction/printdraftjournal"
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        var tabOrWindow = window.open('/res/docs/listjournal/' + res, '_blank');
        tabOrWindow.focus();
    })
}

journal.refresh = function () {
    journal.refreshDraftwithFilter(true)
    journal.getDataJournal(function () {
        journal.renderGrid()
    })
}
journal.refreshDraft = function () {
    journal.refreshDraftwithFilter(true)
    journal.getDataDraft(function () {
        journal.renderGridDraft()
    })
}
journal.setDate = function () {
    var datepicker = $("#dateStartJournal").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    var datestartjournal = $("#dateStartAllJournal").data("kendoDatePicker");
    datestartjournal.value(new Date(newDate))
    journal.dateStart(new Date(newDate))
    datepicker.value(new Date(newDate))
}


journal.getDataDraft = function (callback) {
    var s = $('#dateStartAllJournal').data('kendoDatePicker').value();
    var e = $('#dateEndAllJournal').data('kendoDatePicker').value();
    var param = {
        Type: journal.typeFilterForJournal(),
        DateStart: moment(s).format('YYYY-MM-DD'),
        DateEnd: moment(e).format('YYYY-MM-DD'),
        Filter: journal.refreshDraftwithFilter(),
        TextSearch: journal.textSearch(),
    }
    ajaxPost('/transaction/getdatadraftjournal', param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        journal.dataDraftJournal(res.Data)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

journal.renderGridDraft = function () {
    var data = journal.dataDraftJournal();
    if (typeof $('#listDraft').data('kendoGrid') !== 'undefined') {
        $('#listDraft').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            // pageSize: 25,
        }))
        return
    }

    var columns = [{
        title: 'Delete',
        template: "<a href=\"javascript:journal.deleteDraft('#: _id #')\" class=\"btn btn-danger btn-xs\"><i class=\"fa fa-trash\"></i></a>",
        attributes: {
            "class": "centerAction",
        },
        width: 50,
    }, {
        field: 'IdJournal',
        title: 'No. Journal',
        width: 150,
        template: function (d) {
            return '<a class="onclickJurnal" onclick="journal.renderGridDraftDetail(\'' + d._id + '\',\'' + d.Journal_Type + '\',\'' + d.DateStr + '\',\'' + d.IdJournal + '\',\'' + d.CreateDate + '\')">' + d.IdJournal + '</a>'
        },
        attributes: {
            "class": "colorText",
        }
    }, {
        field: 'DateStr',
        title: 'Date',
        width: 100,
        template: "#=moment(DateStr).format('DD-MMM-YYYY')#"
    }, {
        field: 'Journal_Type',
        title: 'Type',
        width: 100
    }, {
        field: 'User',
        title: 'User',
        width: 60

    },]

    $('#listDraft').kendoGrid({
        dataSource: {
            data: data,
        },
        height: 500,
        sortable: true,
        scrollable: true,
        columns: columns,
    })
    var grid = $("#listDraft").data("kendoGrid");
    if (userinfo.rolenameh() == "admin") {
        grid.hideColumn(0);
    }
}
journal.deleteDraft = function (id) {
    var url = "/transaction/deletedraftjournal"
    var param = {
        Id: id
    }
    swal({
        title: "Are you sure?",
        text: "You will delete this journal",
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
            ajaxPost(url, param, function (res) {
                if (res.IsError == true) {
                    model.Processing(false)
                    return swal({
                        title: "Error!",
                        text: "Error to delete this journal!!!",
                        type: "error",
                        confirmButtonColor: '#3da09a'
                    })
                }
                setTimeout(function () {
                    swal({
                        title: "Success!",
                        text: "Journal has been deleted",
                        type: "success",
                        confirmButtonColor: '#3da09a'
                    });
                    journal.getDataDraft(function () {
                        journal.renderGridDraft()
                    })
                    model.Processing(false)
                }, 100);

            })
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}
journal.DraftAttachments = [];
journal.attachToDraft = function (id) {
    if (!journal.DraftAttachments.includes(id)) {
        journal.DraftAttachments.push(id);
        $("#btnDraftAtt_" + id).addClass("btn-primary");
    }
    var fname = $("#upload_" + id).val();
    var spl = fname.split(/\/|\\/);
    fname = spl[spl.length - 1];
    if (fname.length > 17) {
        fname = fname.substr(0, 7) + "..." + fname.substr(fname.length - 7);
    }
    $("#lbDraftAtt_" + id).html("(" + fname + ")");
}
journal.renderGridDraftDetail = function (ids, journaltype, datestr, nojurnal, createdate) {
    journal.idForPosting(ids);
    journal.createDate(createdate);
    journal.DraftAttachments = [];
    var data = _.filter(journal.dataDraftJournal(), {
        _id: ids
    });
    if (typeof $('#listDraftDetail').data('kendoGrid') !== 'undefined') {
        $('#listDraftDetail').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: (_.get(data, "[0].ListDetail", [])),
            pageSize: 25,
        }))
        if (ids != null) {
            journal.textListDraftDetail(journaltype + "/" + datestr + "/" + nojurnal)
        } else {
            journal.textListDraftDetail("")
        }

        return
    }

    var columns = [
        {
            title: 'Attachment',
            width: 120,
            attributes: {
                "class": "centerAction",
            },
            template: function (e) {
                if (e.Attachment != "") {
                    return "<div class=\"social-box\"><i class=\"glyphicon glyphicon-paperclip white-color\"></i></div> <label> Attached</label>"
                } else {
                    var att = '<input type="file" class="upload" data-bind="attr: {id: \'upload_' + e._id + '\' },value:Attachment" name="upload" id="upload_' + e._id + '">'
                    //return "<a class=\"social-box1 attached\" href=\"javascript:journal.attachToDraft('"+ e._id +"')\"><i class=\"glyphicon glyphicon-paperclip white-color\"></i></a>"
                    //return att + '<label>Not Attached</label>';

                    //return '<input type="file" class="upload" data-bind="attr: {id: \'upload_' + e._id + '\' + $index() },value:Attachment" name="upload" />';

                    return `<div class="upload-btn-wrapper">
                    <button class="btn btn-sm" id="btnDraftAtt_` + e._id + `"><span class="glyphicon glyphicon-upload"></span></button>
                    <input type="file" class="upload" name="upload" id="upload_` + e._id + `" onchange="journal.attachToDraft('` + e._id + `')">
                    <span class="glyphicon glyphicon-paperclip" style="display: none;"></span>
                    </div> <label id="lbDraftAtt_` + e._id + `"> Not Attached</label>`;
                }
            }

        }, {
            field: 'DocumentNumber',
            title: 'Doc No.',
            width: 100
        }, {
            field: 'Acc_Code',
            title: 'Acc Code',
            width: 70
        }, {
            field: 'Acc_Name',
            title: 'Acc Name',
            width: 150
        }, {
            field: 'Department',
            title: 'Department',
            width: 80
        }, {
            title: 'Sales',
            width: 120,
            template: function (d) {
                if (d.SalesName == "" || d.SalesName == undefined) return "--"
                else return d.SalesName
            }
        }, {
            field: 'Debet',
            title: 'Debit',
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            template: "#=ChangeToRupiah(Debet)#"
        }, {
            field: 'Credit',
            title: 'Credit',
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            template: "#=ChangeToRupiah(Credit)#"
        }, {
            field: 'Description',
            title: 'Description',
            width: 100
        }]

    $('#listDraftDetail').kendoGrid({
        dataSource: {
            data: (_.get(data, "[0].ListDetail", [])),
        },
        sortable: true,
        scrollable: true,
        columns: columns,
    })
}

journal.saveToPosting = function () {
    if ($("#listDraftDetail").data("kendoGrid").dataSource.total() == 0) {
        return swal({
            title: "Error!",
            text: "Please Choose List Draft for the first",
            type: "error",
            confirmButtonColor: '#3da09a',
        })
    }

    var date = new Date()
    var dateString = moment(date).format('DD-MM-YYYY')
    var dateNow = moment(dateString, 'DD-MM-YYYY').toDate()
    var cDate = moment(journal.createDate()).format('DD-MM-YYYY')
    var CreatedDate = moment(cDate, "DD-MM-YYYY").toDate()
    var msg = "You will posting this journal"
    if (dateNow > CreatedDate) {
        msg = "You will posting this journal over the created date"
    }
    var url = "/transaction/savepostingex"
    var param = {
        Id: journal.idForPosting(),
        Role: userinfo.rolenameh(),
    }

    // Get Attachments
    var att = {};
    var ids = [];
    for (var k in journal.DraftAttachments) {
        var fi = $("#upload_" + journal.DraftAttachments[k])[0];
        if (fi.files.length > 0) {
            if (!ProActive.checkUploadFilesize(fi.files[0].size,
                "Attachment file size cannot exceed $size!")) return;
            att[journal.DraftAttachments[k]] = fi.files[0];
            ids.push(journal.DraftAttachments[k]);
        }
    }

    // Create Form Data
    var formData = new FormData()
    formData.append('Id', param.Id);
    formData.append('Role', param.Role);
    formData.append('Details', ids);
    for (var k in att) {
        formData.append("fileUpload_" + k, att[k]);
    }

    // Save Data
    swal({
        title: "Are you sure?",
        text: "You will submit this journal",
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
                    if (data.Status != "OK") {
                        model.Processing(false)
                        return swal({
                            title: "Error!",
                            text: "Error Posting Journal!!!",
                            type: "error",
                            confirmButtonColor: '#3da09a'
                        })
                    }
                    swal({
                        title: "Success!",
                        text: "Success Posting Journal",
                        type: "success",
                        confirmButtonColor: '#3da09a'
    
                    });
                    model.Processing(false)
                    journal.init()
                }
            });
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}
journal.postingAndPrint = function () {
    if ($("#listDraftDetail").data("kendoGrid").dataSource.total() == 0) {
        return swal({
            title: "Error!",
            text: "Please Choose List Draft for the first",
            type: "error",
            confirmButtonColor: '#3da09a'
        })
    }

    // Get Attachments
    for (var k in journal.DraftAttachments) {
        var fi = ($("#upload_" + journal.DraftAttachments[k]));
        console.log(fi);
    }

    var date = new Date()
    var dateString = moment(date).format('DD-MM-YYYY')
    var dateNow = moment(dateString, 'DD-MM-YYYY').toDate()
    var cDate = moment(journal.createDate()).format('DD-MM-YYYY')
    var CreatedDate = moment(cDate, "DD-MM-YYYY").toDate()
    var msg = "You will posting this journal"
    if (dateNow > CreatedDate) {
        msg = "You will posting this journal over the created date"
    }
    var url = "/transaction/savepostingandprint"
    var param = {
        Id: journal.idForPosting(),
        Role: userinfo.rolenameh(),

    }
    swal({
        title: "Are you sure?",
        text: msg,
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: false,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, param, function (data) {
                //console.log(data)
                if (data.Message != "OK") {
                    model.Processing(false)
                    return swal({
                        title: "Error!",
                        text: "Error Posting Journal!!!",
                        type: "error",
                        confirmButtonColor: '#3da09a'
                    })
                }
                model.Processing(false)
                swal({
                    title: "Success!",
                    text: "Success Posting Journal",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                }, function () {
                    var tabOrWindow = window.open('/res/docs/journal/' + data.Data, '_blank');
                    tabOrWindow.focus();
                    journal.init()
                });
            })
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}
journal.comeToListofDraft = function () {
    $("#filtered").prop('disabled', false);
    var dropdownList = $('#dropdownTypeDraft').data("kendoDropDownList");
    dropdownList.value("");
    journal.typeFilter("");
    journal.showListDraft(true);
    journal.init();
    document.getElementById('listofdraftjournal').style.display = 'none';
    journal.ShowFilterJournal(true)
    $('#textSearch').val("")
}

journal.refreshall = function () {
    //console.log($('#textSearch').val())
    if (journal.TabsIndexPosition() == 3) {
        journal.textSearch($('#textSearch').val())
        journal.refreshDraft()
    } else if (journal.TabsIndexPosition() == 4) {
        journal.textSearch($('#textSearch').val())
        journal.refresh()
    }
}
journal.exportToExcel = function () {
    console.log("ok")
    model.Processing(true)
    ajaxPost("/transaction/exporttoexcelall", {}, function (res) {
        // var pom = document.createElement('a');
        // pom.setAttribute('href', "/res/docs/journal/" + e.Data);
        // pom.setAttribute('download', e.Data);
        // pom.click();
        var tabOrWindow = window.open('/res/docs/journal/' + res.Data, '_blank');
        tabOrWindow.focus();
        model.Processing(false)
    })
}

journal.checkActiveTab = function (e) {
    if (e == 1 || e == 2) {
        $("#salesOption").show()
    } else {
        $("#salesOption").hide()
    }
    journal.TabsIndexPosition(e)
}
journal.departmentDropdown = function () {
    $("#departmenDropdown").html("")
    var data = []
    ajaxPost("/transaction/getdatadepartment", {}, function (res) {
        $("#departmenDropdown").kendoDropDownList({
            filter: "contains",
            dataTextField: "DepartmentName",
            dataValueField: "DepartmentName",
            dataSource: res.Data,
            optionLabel: 'Select one',
            noDataTemplate: $("#noDataTemplate").html(),
            // open: function() {
            //     $(this.filterInput[0]).val("")
            //     this.dataSource.filter({});
            // },
            change: function (e) {
                var dataitem = this.dataItem();
                journal.valueDepartment(dataitem.DepartmentName)
            }
        });
        var dropdownlist = $("#departmenDropdown").data("kendoDropDownList");
        if (dropdownlist.value() != "") {
            dropdownlist.value("");
            dropdownlist.trigger("change");
        }
    })
}

journal.salesDropdown = function () {
    $("#salesDropdown").html("")
    var data = []
    ajaxPost("/master/getdatasales", {}, function (res) {
        for (i in res.Data) {
            res.Data[i].SalesCode = res.Data[i].SalesID
            res.Data[i].SalesName = res.Data[i].SalesName
        }
        $("#salesDropdown").kendoDropDownList({
            filter: "contains",
            dataTextField: "SalesName",
            dataValueField: "SalesID",
            dataSource: res.Data,
            optionLabel: 'Select one',
            noDataTemplate: $("#noDataTemplate").html(),
            // open: function() {
            //     $(this.filterInput[0]).val("")
            //     this.dataSource.filter({});
            // },
            change: function (e) {
                var dataitem = this.dataItem();
                journal.valueSales(dataitem.SalesID)
                // var res = (dataitem.SalesName).split(" - ");
                journal.nameSales(dataitem.SalesName)
            }
        });
        var dropdownlist = $("#salesDropdown").data("kendoDropDownList");
        if (dropdownlist.value() != "") {
            dropdownlist.value("");
            dropdownlist.trigger("change");
        }
    })
}

journal.addNewDepart = function (widgetId, value, value2) {
    var widget = $("#" + widgetId).getKendoDropDownList();
    //console.log(widgetId, value, value2, widget)    
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
    }, function (isConfirm) {
        if (isConfirm) {
            dataSource.add({
                ID: "",
                DepartmentCode: "",
                DepartmentName: value.toUpperCase()
            });
            dataSource.one("sync", function () {
                widget.select(dataSource.view().length - 1);
            });
            dataSource.sync();
            var dropdownlist = $("#departmenDropdown").data("kendoDropDownList");
            dropdownlist.value(value.toUpperCase());
            dropdownlist.trigger("change");
            model.Processing(true)
            var url = "/transaction/savenewdepartment"
            var param = {
                DepartmentName: value.toUpperCase()
            }
            ajaxPost(url, param, function (data) {
                model.Processing(false)
                swal({
                    title: "Success!",
                    text: "Success Add New Data",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                })
                console.log(data)
            })
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
    // if (confirm("Are you sure?")) {
    //     dataSource.add({
    //         DepartmentCode: 0,
    //         DepartmentName: value
    //     });

    //     dataSource.one("sync", function() {
    //         widget.select(dataSource.view().length - 1);
    //     });

    //     dataSource.sync();
    // }
}
journal.ImportExcel= function () {
    event.stopPropagation();
    event.preventDefault();

    model.Processing(true);

    if ($('input[type=file]')[0].files[0] == undefined) {
        swal('Error', 'Please select a file to upload!', 'error');
        model.Processing(false);
        return;
    }

    if ($('input[type=file]')[0].files[0].size > 1000000) {
        swal('Out of Limit', 'Allowed file size is 1 Mb!', 'error');
        model.Processing(false);
        return;
    }
    
    var len = $("#fDok")[0].files.length
    if (len > 0) {
        var j = 0;
        var data = new FormData();
        for (i = 0; i < $("#fDok")[0].files.length; i++) {
            data.append("filedoc", $('input[type=file]')[0].files[i]);
            data.append("filename", $('input[type=file]')[0].files[i].name);
        }
    } else {
        var data = new FormData();
        data.append("filedoc", $('input[type=file]')[0].files[0]);
        data.append("filename", $('input[type=file]')[0].files[0].name);
    }
    if ($('input[type=file]')[0].files[0].name != "") {
        jQuery.ajax({
            url: '/transaction/uploadfiles',
            data: data,
            cache: false,
            contentType: false,
            processData: false,
            type: 'POST',
            success: function (data) {
                swal('Success', 'Data has been uploaded successfully!', 'success');
                model.Processing(false);
                $('#fDok').val('');
                $('#fInfo1').val('');
                $('#ImportAllJournalModal').modal('hide');
                journal.init()
            }
        });
    } else {
        swal('Error', 'Please select a file to upload!', 'error');
        model.Processing(false);
    }
}

model.isFormValid = function (selector) {
    model.resetValidation(selector);
    var $validator = $(selector).data("kendoValidator");
    return ($validator.validate());
};

model.resetValidation = function (selectorID) {
    var $form = $(selectorID).data("kendoValidator");
    if ($form == undefined) {
        $(selectorID).kendoValidator();
        $form = $(selectorID).data("kendoValidator");
    }

    $form.hideMessages();
};
journal.fromOtherPage = function () {
    var url_string = window.location.href
    var url = new URL(url_string);
    var num = url.searchParams.get("id");
    if (num != null) {
        if (num != "multijournal") {
            journal.idJournalFromCoa(num)
            journal.idJournalFromCoaForUnapply(num)
            journal.showUnApply(true)
            journal.refresh()
        } else {
            //console.log("tes")
            setTimeout(function () {
                $('.nav-tabs a[href="#listofdraft"]').tab('show');
                journal.comeToListofDraft()
                journal.checkActiveTab(3)
            }, 100)

        }
        history.pushState(null, null, '/transaction/journal');
    }
}
journal.openModalUnapply = function () {
    if (journal.statusJournal() == "close") {
        return swal({
            title: "Warning!",
            text: "This journal was closed!!!",
            type: "warning",
            confirmButtonColor: '#3da09a'
        })
    } else {
        $("#unapplyJournal").modal("show");
    }
}
journal.unAplly = function () {
    var date = $('#DateUnapply').data('kendoDatePicker').value()
    var journalNumber = journal.idJournalFromCoaForUnapply()
    var url = "/transaction/unapplyjournal"
    var param = {
        Id: journalNumber,
        Date: moment(date).format('YYYY-MM-DD'),
    }
    swal({
        title: "Are you sure?",
        text: "You will save unapply data with this date",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: false,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, param, function (data) {
                //console.log(data)
                if (data.IsError != false) {
                    model.Processing(false)
                    return swal({
                        title: "Error!",
                        text: "Error Posting Journal!!!",
                        type: "error",
                        confirmButtonColor: '#3da09a'
                    })
                }
                model.Processing(false)
                swal({
                    title: "Success!",
                    text: "Success Posting Journal",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                }, function () {
                    $("#unapplyJournal").modal("hide");
                    journal.init()
                });
            })
        } else {
            swal({
                title: "Cancelled",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
    });
}
journal.init = function () {
    $(document).on('change', ':file', function () {
        // ValidateFiles($(this));
        var input = $(this),
            numFiles = input.get(0).files ? input.get(0).files.length : 1,
            label = input.val().replace(/\\/g, '/').replace(/.*\//, '');
        input.trigger('fileselect', [
            numFiles, label]);
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
    journal.getDateNow()
    journal.setDate()
    journal.getData();
    journal.cashNameJournal("Journal")
    journal.titleListJournal("All")
    journal.getDataJournal(function () {
        journal.renderGrid()
    })
    journal.getDataDraft(function () {
        journal.renderGridDraft()
    })
    journal.renderGridDraftDetail()
}
journal.onChangeDateStart = function (val) {
    if (val.getTime() > journal.dateEnd().getTime()) {
        journal.dateEnd(val)
    }
}

journal.filterText = function (term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["Acc_Name", "DocumentNumber", "Department", "SalesName", "Description"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator: "contains", value: $searchValue });
    }
    $("#listJournal").data("kendoGrid").dataSource.query({ filter: filter, pageSize: 25, page: 1 });

}

$(document).ready(function () {
    journal.checkActiveTab(4)
    journal.fromOtherPage()
    $('.nav-tabs a[href="#listofjournal"]').tab('show')
    journal.comeToListJournal()
    journal.init();
    journal.record.ListDetail([]);
    journal.getLastDocNumber()
    $(".btn").mouseup(function () {
        $(this).blur();
    })
    $("#textSearch").on("keyup blur change", function () {
        journal.filterText();
    });

})