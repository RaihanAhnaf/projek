model.Processing(false)
var glptcmandiri = {}
glptcmandiri.koFilterIsActive = ko.observable(false)
glptcmandiri.koDataPettyCash = ko.observableArray([])
glptcmandiri.refreshGrid = function () {
    glptcmandiri.koFilterIsActive(true)
    glptcmandiri.getDataglptcmandiri(function () {
        glptcmandiri.renderGrid()
    })
}
glptcmandiri.getDataglptcmandiri = function (callback) {
    model.Processing(true)
    var url = "/report/getdataglptcmandiri"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: glptcmandiri.koFilterIsActive(),
    }
    ajaxPost(url, param, function (res) {
        glptcmandiri.koDataPettyCash(res.Data)
        model.Processing(false)
        callback()
    })
}
glptcmandiri.renderGrid = function () {
    var data = glptcmandiri.koDataPettyCash()
    var columns = [{
        field: 'PostingDate',
        title: 'Posting Date',
        template: function (e) {
            return moment.utc(e.PostingDate).format("DD MMM YYYY")
        },
        width: 70,
    }, {
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 100
    }, {
        field: 'Description',
        title: 'Description',
        width: 200,
        footerTemplate: "Total :"
    }, {
        field: 'Debet',
        title: 'Debit',
        width: 100,
        template: "#=ChangeToRupiah(Debet)#",
        footerTemplate: "#= ChangeToRupiah(sum) #"
    }, {
        field: 'Credit',
        title: 'Credit',
        width: 100,
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: "#= ChangeToRupiah(sum) #"
    }, {
        field: 'Saldo',
        title: 'Saldo',
        width: 100,
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: function (e) {
            return glptcmandiri.saldoCalculation()
        }
    }, ]
    $('#gridGLPTCMandiri').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Saldo",
                aggregate: "sum"
            }],
        },
        height: 500,
        filterable: true,
        scrollable: true,
        columns: columns
    })
}
glptcmandiri.saldoCalculation = function () {
    var data = $('#gridGLPTCMandiri').data('kendoGrid').dataSource.options.data
    var sumDebit = _.sumBy(data, 'Debet')
    var sumCredt = _.sumBy(data, 'Credit')
    var Saldo = sumDebit - sumCredt
    return kendo.toString(Saldo, 'n')
}
glptcmandiri.init = function () {
    glptcmandiri.getDataglptcmandiri(function () {
        glptcmandiri.renderGrid()
    })
}
$(document).ready(function () {
    glptcmandiri.init()
})