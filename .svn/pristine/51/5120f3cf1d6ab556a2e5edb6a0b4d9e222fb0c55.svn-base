model.Processing(false)
var multijournal = {}

multijournal.docnofromdb = ko.observable()
multijournal.showListJournal = ko.observable(false)
multijournal.showListDraft = ko.observable(false)
multijournal.PrintiIsActive = ko.observable(false)
multijournal.dataJurnal = ko.observableArray([])
multijournal.cashNameJournal = ko.observable("")
multijournal.balanceData = ko.observable(false)
multijournal.documentNoStr = ko.observable("")
multijournal.dataMasterAccount = ko.observableArray([])
multijournal.accountNumberForDropDown = ko.observableArray([])
multijournal.documentNumberLastRow = ko.observable("")
multijournal.textListDraftDetail = ko.observable("")
multijournal.typeFilter = ko.observable("")
multijournal.titleListJournal = ko.observable("")
multijournal.dataMasterJournal = ko.observableArray([])
multijournal.dataDraftJournal = ko.observableArray([])
multijournal.dateStartFilter = ko.observable()
multijournal.dateEndFilter = ko.observable("")
multijournal.showIconAttachment = ko.observable(false)
multijournal.showPdf = ko.observable(false)
multijournal.typeFilterForJournal = ko.observable('')
multijournal.idForPosting = ko.observable()
multijournal.createDate = ko.observable("")
multijournal.refreshDraftwithFilter = ko.observable(false)
multijournal.lastDate = ko.observable()
multijournal.DatePageBar = ko.observable()
multijournal.ActiveTab = ko.observable(0)
multijournal.dataDepartment= ko.observableArray([])
multijournal.dataSales= ko.observableArray([])
multijournal.typeAllJournal = [{
    "value": "CashIn",
    "text": "Cash In"
}, {
    "value": "CashOut",
    "text": "Cash Out"
}, {
    "value": "General",
    "text": "General"
}]

multijournal.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    multijournal.DatePageBar(page)
}

multijournal.cashNameJournal.subscribe(function (e) {
    $('#textSearch').val("")
})

multijournal.newRecord = function () {
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
        Department:"MULTIJOURNAL",
        SalesCode: "",
        SalesName: ""
    }
    page.ListDetail.push(multijournal.listDetailJournal())
    return page
}

multijournal.listDetailJournal = function () {
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
        Department: ko.observable(''),
        SalesCode: ko.observable(''),
        SalesName: ko.observable(''),
        Attachment: ko.observable(''),
        User: ko.observable(''),
    }
}

multijournal.record = ko.mapping.fromJS(multijournal.newRecord())

multijournal.getData = function () {
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
        multijournal.dataMasterAccount(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].ACC_Code = DataAccount[i].ACC_Code + ""
            DataAccount[i].CodeName = DataAccount[i].ACC_Code + "-" + DataAccount[i].Account_Name
        }
        multijournal.accountNumberForDropDown(DataAccount)

        model.Processing(false)
    })
}

multijournal.totalDebet = ko.computed(function () {
    if (!_.isArray(multijournal.record.ListDetail())) return 0;
    var totalDebet = _.sumBy(multijournal.record.ListDetail(), function (v) {
        var debet = kendo.toString(v.Debet(), "n");
        var debetInt = Number(debet.replace(/[^0-9\.]+/g, ""));
        return debetInt
    })
    return kendo.toString(totalDebet, "n");
})

multijournal.totalCredit = ko.computed(function () {
    if (!_.isArray(multijournal.record.ListDetail())) return 0;
    var totalCredit = _.sumBy(multijournal.record.ListDetail(), function (v) {
        var credit = kendo.toString(v.Credit(), "n");
        var creditInt = Number(credit.replace(/[^0-9\.]+/g, ""));
        return creditInt
    })
    return kendo.toString(totalCredit, "n");
})

multijournal.balance = ko.computed(function () {
    var totalDebet = Number(multijournal.totalDebet().replace(/[^0-9\.]+/g, ""));
    var totalCredit = Number(multijournal.totalCredit().replace(/[^0-9\.]+/g, ""));
    var total = totalDebet - totalCredit

    if (total >= 0) {
        var TotString = kendo.toString(total, "n");
        return TotString;
    } else {
        var TotminString = kendo.toString(Math.abs(total), "n");
        return "(" + TotminString + ")";
    }

})
multijournal.onChangeDate = function (value, index){
    // console.log(value, index,  multijournal.record.ListDetail()[index].DocumentNumber())
    var formarDate = moment(value).format("DDMMYY")
    var split = multijournal.record.ListDetail()[index].DocumentNumber().split("/")
    var newDoc = split[0]+"/"+formarDate +"/"+split[2]+"/"+split[3]
    multijournal.record.ListDetail()[index].DocumentNumber(newDoc);
}
multijournal.onChangeAccountNumber = function (value, index) {
    findaccount = _.find(multijournal.dataMasterAccount(), {
        ACC_Code: value
    })
    multijournal.record.ListDetail()[index].Acc_Name(findaccount.Account_Name);
}

multijournal.addContent = function () {
    var documentcode = ""
    var journaltype = ""
    multijournal.record.ListDetail.push(multijournal.listDetailJournal())
    multijournal.record.User(userinfo.usernameh());
    multijournal.record.CreateDate(new Date());
    multijournal.record.PostingDate(new Date());
    multijournal.record.DateStr(moment().format("DD MMMM YYYY"));
    multijournal.record.Status("draft");
    multijournal.record.ListDetail()[multijournal.record.ListDetail().length - 1].No(multijournal.record.ListDetail().length);
    multijournal.record.ListDetail()[multijournal.record.ListDetail().length - 1].PostingDate(new Date());
    multijournal.record.ListDetail()[multijournal.record.ListDetail().length - 1].User(userinfo.usernameh());

    switch (multijournal.cashNameJournal()) {
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
    if (multijournal.record.ListDetail().length < 2) {

        number = 1

    } else {

        if (multijournal.balance() == "0.00" && multijournal.totalDebet() != "0.00") {
            var jr = ko.mapping.toJS(multijournal.record)
            var i = multijournal.record.ListDetail().length - 2
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
        } else if (multijournal.totalDebet() == "0.00" || multijournal.totalCredit() == "0.00") {
            number = multijournal.documentNumberLastRow()
        } else {
            number = multijournal.documentNumberLastRow()
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
    multijournal.record.Journal_Type(journaltype);
    multijournal.record.ListDetail()[multijournal.record.ListDetail().length - 1].Id(multijournal.makeId());
    multijournal.record.ListDetail()[multijournal.record.ListDetail().length - 1].DocumentNumber(documentcode + "/" + moment().format("DDMMYY") + "/temp-" + userinfo.usernameh() + "/" + beforeNumber + number);
    multijournal.record.ListDetail()[multijournal.record.ListDetail().length - 1].Journal_Type(journaltype);
    multijournal.maskingMoney()
    multijournal.documentNumberLastRow(number)
}

multijournal.getLastDocNumber = function () {
    switch (multijournal.cashNameJournal()) {
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
            multijournal.docnofromdb(parseInt(no))
        }

    })
}

multijournal.addNewRow = function () {
    if (multijournal.cashNameJournal() == "") {
        swal({
            title:"Error!",
            text: "Please Choose Journal Type for the first",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
        return
    }
    if (multijournal.record.ListDetail().length > 0) {
        var jr = ko.mapping.toJS(multijournal.record)
        var i = multijournal.record.ListDetail().length - 1
        var no = jr.ListDetail[i].DocumentNumber.substr(jr.ListDetail[i].DocumentNumber.length - 1);
        var index = multijournal.record.ListDetail().length;
        multijournal.addContent()
        if (multijournal.ActiveTab() == 0) {
            setTimeout(function () {
                $("#sales"+ index).data("kendoDropDownList").enable(false); 
            }, 100);
        }
        return
    }
    multijournal.addContent()
    if (multijournal.ActiveTab() == 0) {
        // console.log("hhher")
        setTimeout(function () {
            $("#sales0").data("kendoDropDownList").enable(false); 
        }, 100);
    }
}

multijournal.maskingMoney = function () {
    $('.currency').inputmask("numeric", {
        radixPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        rightAlign: false,
    });
}

multijournal.removeRow = function () {
    multijournal.record.ListDetail.remove(this)
    if (multijournal.record.ListDetail().length != 0) {
        var jr = ko.mapping.toJS(multijournal.record)
        var i = multijournal.record.ListDetail().length - 1
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
        multijournal.documentNumberLastRow(Nom)
    }
}

multijournal.changetoCashIn = function () {
    multijournal.ActiveTab(0)
    if (multijournal.record.ListDetail().length> 1) {
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
        }, function (isConfirm) {
            if (isConfirm) {
                multijournal.record.ListDetail([]);
                multijournal.cashNameJournal("Cash In")
                multijournal.addNewRow()
            }else{
                if (multijournal.cashNameJournal()=="Cash Out"){
                    $(".tab1").removeClass("active")
                    $(".tab2").addClass("active")
                }
                if (multijournal.cashNameJournal()=="General Journal"){
                    $(".tab1").removeClass("active")
                    $(".tab3").addClass("active")
                }
            }
        });
    } else {
        multijournal.record.ListDetail([]);
        multijournal.cashNameJournal("Cash In")
        multijournal.addNewRow()
    }
}

multijournal.changetoCashOut = function () {
    multijournal.ActiveTab(1)
    if (multijournal.record.ListDetail().length > 1) {
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
        },function (isConfirm) {
            if (isConfirm) {
                multijournal.record.ListDetail([]);
                multijournal.cashNameJournal("Cash Out")
                multijournal.addNewRow()
            }else{
                if (multijournal.cashNameJournal()=="Cash In"){
                    $(".tab2").removeClass("active")
                    $(".tab1").addClass("active")
                }
                if (multijournal.cashNameJournal()=="General Journal"){
                    $(".tab2").removeClass("active")
                    $(".tab3").addClass("active")
                }
            }
        });
    } else {
        multijournal.record.ListDetail([]);
        multijournal.cashNameJournal("Cash Out")
        multijournal.addNewRow()
    }

}

multijournal.changetoGeneral = function (type) {
    multijournal.ActiveTab(2)
    if (multijournal.record.ListDetail().length > 1) {
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
        }, function (isConfirm) {
            if (isConfirm) {
                multijournal.record.ListDetail([]);
                multijournal.cashNameJournal("General Journal")
                multijournal.addNewRow()
            }else{
                if (multijournal.cashNameJournal()=="Cash In"){
                    $(".tab3").removeClass("active")
                    $(".tab1").addClass("active")
                }
                if (multijournal.cashNameJournal()=="Cash Out"){
                    $(".tab3").removeClass("active")
                    $(".tab2").addClass("active")
                }
            }
        });
    } else {
        multijournal.record.ListDetail([]);
        multijournal.cashNameJournal("General Journal")
        multijournal.addNewRow()
    }

}

multijournal.makeId = function () {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 10; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}
multijournal.resetRow = function () {
    multijournal.record.ListDetail([]);
}

multijournal.saveData = function () {
    var payload = ko.mapping.toJS(multijournal.record)
    var formData = new FormData()
    if (payload.ListDetail.length == 0) {
        return swal({
            title:"Warning!",
            text: "None data for save",
            type: "info",
            confirmButtonColor:"#3da09a"})
    }
    var i = multijournal.record.ListDetail().length - 1
    var accname = payload.ListDetail[i].Acc_Code
    if (accname == "") {
        return swal({
            title:"Warning!", 
            text:"Please check Account Code",
            type: "info",
            confirmButtonColor:"#3da09a"})
    } else if (multijournal.totalDebet() <= 0 || multijournal.totalCredit() <= 0 || multijournal.balance() != 0) {
        return swal({
            title:"Warning!",
            text: "Please Balance the Journal for the first", 
            type:"info",
            confirmButtonColor:"#3da09a"})
    }

    var indexList = 0
    var fillDept = false
    var fillSales = false
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
        
        if (multijournal.ActiveTab()!= 0) {
            // var res = ($("#sales"+i).data("kendoDropDownList").text()).split(" - ");
            payload.ListDetail[i].SalesName = $("#sales"+i).data("kendoDropDownList").text()
        }
    })
    payload.Department = "MULTIJOURNAL"
    payload.SalesCode = "MULTIJOURNAL"
    payload.SalesName = "MULTIJOURNAL"
    formData.append('data', JSON.stringify(payload))
    var i
    var attachment = document.getElementsByClassName('upload');
    for (i = 0; i < attachment.length; i++) {
        formData.append("fileUpload" + i, attachment[i].files[0]);
    }
    for (i in payload.ListDetail) {
        if (payload.ListDetail[i].Department == "") {
            indexList= i
            fillDept = true
            break
        }
    }
    if(fillDept){
        return swal({
            title:"Warning!",
            text: "Please check department in row "+indexList, 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }

    if (multijournal.ActiveTab()!= 0) {
        for (i in payload.ListDetail) {
            if (payload.ListDetail[i].SalesCode == "") {
                indexList= i
                fillSales = true
                break
            }
        }
        if(fillSales){
            return swal({
                title:"Warning!",
                text: "Please check sales in row "+indexList, 
                type:"info",
                confirmButtonColor:"#3da09a"
            })
        }
    }
    
    var url = "/transaction/savemultijournal"
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
                    setTimeout(function(){
                        swal({
                            title: "Success!",
                            text: "Journal has been saved!",
                            type: "success",
                            confirmButtonColor:"#3da09a"
                        }, function () {
                            window.location.assign("/transaction/journal?id=multijournal")
                        });
                    },1000)
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
multijournal.getDepartment = function(){
    // $("#departmenDropdown").html("")
    var data = []
    ajaxPost("/transaction/getdatadepartment", {}, function(res){
        multijournal.dataDepartment(res.Data)
    })
}

multijournal.getSales = function(){
    var data = []
    ajaxPost("/master/getdatasales", {}, function(res){
        for (i in res.Data) {
            res.Data[i].SalesCode = res.Data[i].SalesID
            res.Data[i].SalesName = res.Data[i].SalesName
        }
        multijournal.dataSales(res.Data)
    })
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
multijournal.addNewDepart = function(widgetId, value, value2){
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
                Id: "",
                DepartmentCode: "",
                DepartmentName: value.toUpperCase()
            });
            dataSource.one("sync", function() {
                widget.select(dataSource.view().length - 1);
            });
            dataSource.sync();
            var dropdownlist = $("#"+ widgetId).data("kendoDropDownList");
            dropdownlist.value(value.toUpperCase());
            dropdownlist.trigger("change");
            model.Processing(true)
            var url = "/transaction/savenewdepartment"
            var param ={
                DepartmentName : value.toUpperCase()
            }
            ajaxPost(url, param, function (data) {
                // if (data.Message != "OK") {
                //     model.Processing(false)
                //     return swal({
                //         title:"Error!",
                //         text: "Error Posting Journal!!!", 
                //         type:"error",
                //         confirmButtonColor: '#3da09a'
                //     })
                // }
                model.Processing(false)
                swal({
                    title:"Success!",
                    text: "Success Add New Data",
                    type: "success",
                    confirmButtonColor: "#3da09a"
                })
            })
        } else {
            swal({
                title:"Cancelled",
                type:"error",
                confirmButtonColor:"#3da09a"
            });
        }
    });
}
multijournal.init = function () {
    multijournal.getDepartment()
    multijournal.getSales()
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
    multijournal.getData();
    multijournal.cashNameJournal("Cash In")
    multijournal.getDateNow()
}

$(function () {
    multijournal.init()
    multijournal.record.ListDetail([]);
    multijournal.getLastDocNumber()
    $(".btn").mouseup(function () {
        $(this).blur();
    })
    multijournal.addNewRow()
})