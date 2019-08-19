model.Processing(false)
var balset = {}
balset.dataAktiva = ko.observableArray([])
balset.dataOriAktiva = ko.observableArray([])
balset.dataPassiva = ko.observableArray([])
balset.dataOriPassiva = ko.observableArray([])
balset.visibleEditAktiva = ko.observable(true)
balset.visibleEditPassiva = ko.observable(true)
balset.DatePageBar = ko.observable()

balset.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    balset.DatePageBar(page)
}

balset.getDataForAktiva = function (callback) {
    var url = "/balancesheetsetting/getlistaktiva"
    var param = {
        category: "BALANCE SHEET"
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        balset.dataOriAktiva($.extend(true, [], res.Data))
        balset.dataAktiva($.extend(true, [], res.Data))
        model.Processing(false)
        callback();
    });
}
balset.renderGridAktiva = function () {
    $("#listOfAktiva").kendoTreeList({
        dataSource: {
            data :balset.dataAktiva(),
            sort: [{
                field: "ACC_Code",
                dir: "asc"
            }],
        },  
        height: 500,
        columns: [{
            field: "ACC_Code",
            title: "Account Code",
            width: 120
        }, {
            field: "Account_Name",
            title: "Account Name"
        }, {
            field: "Type",
            title: "Type",
            width: 80
        }, {
            title: "Active",
            width: 70,
            attributes: {
                "class": "align-center"
            },
            template: "<input id='acktive#:ID#' class='checkedActiva' onclick='balset.changeActivaActive(\"#:ID#\")' type='checkbox' #: Active==true ? 'checked' : ''# disabled/>"
        }, ]
    })
}
balset.changeActivaActive = function (ID) {
    var data = balset.dataAktiva()
    for (i in data) {
        if (data[i].ID == ID) {
            var check = $("#acktive" + ID).is(':checked')
            if (check == true) {
                data[i].Active = true
            } else {
                data[i].Active = false
            }
        }
    }
    balset.dataAktiva(data)
}
balset.editAktiva = function () {
    balset.visibleEditAktiva(false)
    $(".checkedActiva").prop('disabled', false);
}
balset.cancelEditAktiva = function () {
    balset.visibleEditAktiva(true)
    $(".checkedActiva").prop('disabled', true);
    $("#listOfAktiva").data("kendoTreeList").setDataSource(balset.dataOriAktiva());
    $("#listOfAktiva").data('kendoTreeList').refresh();
}
balset.getDataForPassiva = function (callback) {
    var url = "/balancesheetsetting/getlistpassiva"
    var param = {
        category: "BALANCE SHEET"
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        balset.dataPassiva($.extend(true, [], res.Data))
        balset.dataOriPassiva($.extend(true, [], res.Data))
        model.Processing(false)
        callback();
    });
}
balset.renderGridPassiva = function () {
    $("#listOfPassiva").kendoTreeList({
        dataSource: {
            data :balset.dataPassiva(),
            sort: [{
                field: "ACC_Code",
                dir: "asc"
            }],
        },
        height: 500,
        columns: [{
            field: "ACC_Code",
            title: "Account Code",
            width: 120
        }, {
            field: "Account_Name",
            title: "Account Name"
        }, {
            field: "Type",
            title: "Type",
            width: 80
        }, {
            title: "Active",
            width: 70,
            attributes: {
                "class": "align-center"
            },
            template: "<input id='acktive#:ID#' class='checkedPassiva' onclick='balset.changeActivaPassiva(\"#:ID#\")' type='checkbox' #: Active==true ? 'checked' : ''# disabled/>"
        }, ]
    })
}
balset.changeActivaPassiva = function (ID) {
    var data = balset.dataPassiva()
    for (i in data) {
        if (data[i].ID == ID) {
            var check = $("#acktive" + ID).is(':checked')
            if (check == true) {
                data[i].Active = true
            } else {
                data[i].Active = false
            }
        }
    }
    balset.dataPassiva(data)
}
balset.editPassiva = function () {
    balset.visibleEditPassiva(false)
    $(".checkedPassiva").prop('disabled', false);
}
balset.cancelEditPassiva = function () {
    balset.visibleEditPassiva(true)
    $(".checkedPassiva").prop('disabled', true);
    $("#listOfPassiva").data("kendoTreeList").setDataSource(balset.dataOriPassiva());
    $("#listOfPassiva").data('kendoTreeList').refresh();
}
balset.SaveDataAktiva = function () {
    var url = "/balancesheetsetting/savedataaccountblaktiva"
    var param = {
        Data: balset.dataAktiva()
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        balset.visibleEditAktiva(true)
        $(".checkedActiva").prop('disabled', true);
        swal({
            title:"Success!",
            text: "Data aktiva for balance sheet is changed",
            type: "success",
            confirmButtonColor:"#3da09a"});
        setTimeout(function () {
            location.reload();
        }, 2000);
    })
}
balset.SaveDataPassiva = function () {
    var url = "/balancesheetsetting/savedataaccountblpassiva"
    var param = {
        Data: balset.dataPassiva()
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        balset.visibleEditAktiva(true)
        $(".checkedPassiva").prop('disabled', true);
        swal({
            title:"Success!",
            text: "Data passiva for balance sheet is changed",
            type: "success",
            confirmButtonColor:"#3da09a"});
        setTimeout(function () {
            location.reload();
        }, 2000);
    })
}
balset.init = function () {
    balset.getDataForAktiva(function () {
        balset.renderGridAktiva()
    })
    balset.getDataForPassiva(function () {
        balset.renderGridPassiva()
    })
    balset.getDateNow()
}
$(document).ready(function () {
    balset.init()
})