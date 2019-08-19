model.Processing(false)
var glpettycash = {}
glpettycash.koFilterIsActive = ko.observable(false)
glpettycash.koDataPettyCash = ko.observableArray([])
glpettycash.refreshGrid = function () {
    glpettycash.koFilterIsActive(true)
    glpettycash.getDataGLpettyCash(function () {
        glpettycash.renderGrid()
    })
}
glpettycash.getDataGLpettyCash = function (callback) {
    model.Processing(true)
    var url = "/report/getdataglpettycash"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: glpettycash.koFilterIsActive(),
    }
    ajaxPost(url, param, function (res) {
        glpettycash.koDataPettyCash(res.Data)
        model.Processing(false)
        callback()
    })
}
glpettycash.renderGrid = function () {
    var data = glpettycash.koDataPettyCash()
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
            return glpettycash.saldoCalculation()
        }
    }, ]
    $('#gridGLpettyCash').kendoGrid({
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
glpettycash.saldoCalculation = function () {
    var data = $('#gridGLpettyCash').data('kendoGrid').dataSource.options.data
    var sumDebit = _.sumBy(data, 'Debet')
    var sumCredt = _.sumBy(data, 'Credit')
    var Saldo = sumDebit - sumCredt
    return kendo.toString(Saldo, 'n')
}
glpettycash.init = function () {
    glpettycash.getDataGLpettyCash(function () {
        glpettycash.renderGrid()
    })
}
$(document).ready(function () {
    glpettycash.init()
})