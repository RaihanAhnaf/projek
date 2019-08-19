model.Processing(false)
var mCoa = {
    DateFilter: ko.observable(false),
    dateNow: ko.observable(),
    MonthNow: ko.observable(),
    YearNow: ko.observable(),
    Period: ko.observable(),
    ShowHideFilter: ko.observable(false),
    categoryCOA: ko.observable(),
    DebetCredit: ko.observable(),
    Datadebetcredit: ko.observableArray([]),
    DataCategory: ko.observableArray([]),
    DataACC_Code: ko.observableArray([]),
    showDetailCoa: ko.observable(false),
    HiddenFilterCoa: ko.observable(false),
    detailCoaName: ko.observable(),
    detailCoaCode: ko.observable(),
    detailCoaParent: ko.observable(),
    DatePageBar: ko.observable(),
    textSearch : ko.observable(""),
    dataDetailCoa : ko.observableArray([]),
    dataMasterDetailCoa : ko.observableArray([]),
    dateStart : ko.observable(new Date),
    dateEnd : ko.observable(new Date)
}
var Datadebetcredit = [{
    "text": "DEBIT",
    "value": "DEBET",
    "background": "greyBackground"
}, {
    "text": "CREDIT",
    "value": "CREDIT"
}]
mCoa.Datadebetcredit(Datadebetcredit)
var DataCategory = [{
    "text": "BALANCE SHEET",
    "value": "BALANCE SHEET"
}, {
    "text": "INCOME STATEMENT",
    "value": "INCOME STATEMENT"
}]

mCoa.getDateNow = function () {
        var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
        mCoa.DatePageBar(page)
    }

mCoa.DataCategory(DataCategory)
var coa = {
    DateShowing: function () {
        var now = new Date()
        mCoa.Period(moment(now).format("DD MMMM YYYY"))
        mCoa.dateNow(moment(now).date())
        mCoa.MonthNow(moment(now).format("MMMM"))
        mCoa.YearNow(moment(now).year())
    },
    GetDataCoa: function () {
        model.Processing(true)
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
        var url = "/master/getdatacoa"
        var param = {}
        if (mCoa.DateFilter() == true) {
            param = {
                Start: dateStart,
                End: dateEnd,
                Filter: true,
                TextSearch: mCoa.textSearch().toUpperCase(),
            }
            // if (param.TextSearch == "") {
            //     param.Filter = false
            // }
        } else {
            param = {
                Filter: false
            }
        }
        ajaxPost(url, param, function (res) {

            model.Processing(false)
            if (res.IsError) {
            swal("Search Not Found!", res.Message, "warning")
            $('#textSearch').val("")
            return
            }
            for (i in res.Data) {
                res.Data[i].ACC_Code = res.Data[i].ACC_Code + ""
            }
            var Data = _.map(res.Data, 'ACC_Code')
            console.log("data => ",res.Data)
            mCoa.DataACC_Code(Data)
            $('#gridMasterCOA').html("");
            $('#gridMasterCOA').kendoGrid({
                dataSource: {
                    data: res.Data,
                    schema: {
                        model: {
                            field: {
                                ACC_Code: {
                                    type: "string"
                                },
                                Account_Name: {
                                    type: "string"
                                },
                                Debet: {
                                    type: "number"
                                },
                                Credit: {
                                    type: "number"
                                },
                                Saldo: {
                                    type: "number"
                                },
                                Debet_Credit: {
                                    type: "string"
                                },
                                Category: {
                                    type: "string"
                                }
                            }
                        }
                    },

                    sort: {
                        field: 'ACC_Code',
                        dir: 'asc',
                    }
                },
                excelExport: function (e) {
                    var rows = e.workbook.sheets[0].rows;
                    for (var ri = 0; ri < rows.length; ri++) {
                        var row = rows[ri];
                        if (row.type == "group-footer" || row.type == "footer") {
                            for (var ci = 0; ci < row.cells.length; ci++) {
                                var cell = row.cells[ci];
                                if (cell.value) {
                                    // Use jQuery.fn.text to remove the HTML and get only the text
                                    var text = $(cell.value).text();
                                    if (cell.value != "Total") {
                                        var expression = /^\([\d,\.]*\)$/;
                                        if (text.match(expression)) {
                                            //It matched - strip out parentheses and append - at front
                                            var val = '-' + text.replace(/[\(\)]/g, '')
                                        } else {
                                            var val = text
                                        }
                                        cell.value = parseFloat(val.split(",").join("")) // this different
                                        cell.format = "#,##0.00_);(#,##0.00);0.00;"
                                        // Set the alignment
                                        cell.hAlign = "right";
                                    } else {
                                        cell.hAlign = "right";
                                    }
                                }
                            }
                        }
                        if (row.type == "header") {
                            for (var ci = 0; ci < row.cells.length; ci++) {
                                var cell = row.cells[ci];
                                cell.value = cell.value.replace("<br/>", "")
                            }
                        }
                        if (row.type == "data") {
                            if (ri > 0) {
                                for (var ci = 0; ci < row.cells.length; ci++) {
                                    var cell = row.cells[ci];
                                    if (ci == 2 || ci == 3 || ci == 4) {
                                        cell.format = "#,##0.00_);(#,##0.00);0.00;"
                                        // Set the alignment
                                        cell.hAlign = "right";
                                    }
                                }
                            }
                        }
                    }
                    e.workbook.fileName = "COA"
                },
                scrollable: true,
                sortable: true,
                height: 500,

                columns: [{
                    field: "ACC_Code",
                    title: "Account Code",
                    width: 120,

                }, {
                    field: "Account_Name",
                    title: "Account Name",
                    template: "<a class='onclickCOA' href='javascript:void(0)' onclick='coa.showDetail(\"#: ID #\",\"#: Account_Name #\",\"#: ACC_Code #\", \"#: Main_Acc_Code #\")'>#:Account_Name#</a>",
                    width: 270,
                    filterable: false
                }, {
                    field: "Debet",
                    title: "Debit",
                    width: 100,
                    filterable: false,
                    template: function (e) {
                        return kendo.toString(e.Debet, "n");
                    },
                    attributes: {
                        style: "text-align:right;"
                    }

                }, {
                    field: "Credit",
                    title: "Credit",
                    width: 100,
                    filterable: false,
                    template: function (e) {
                        return kendo.toString(e.Credit, "n");
                    },
                    attributes: {
                        style: "text-align:right;"
                    }
                }, {
                    field: "Saldo",
                    title: "Saldo",
                    width: 100,
                    filterable: false,
                    template: function (e) {
                        if (e.Saldo >= 0) {
                            var TotString = kendo.toString(e.Saldo, "n");
                            return TotString;
                        } else {
                            var TotminString = kendo.toString(Math.abs(e.Saldo), "n");
                            return "(" + TotminString + ")";
                        }
                    },
                    attributes: {
                        style: "text-align:right;"
                    }
                }, {
                    field: "Debet_Credit",
                    title: "Debit/Credit",
                    width: 100,
                    filterable: false

                }, {
                    field: "Category",
                    title: "Category",
                    width: 120,
                    filterable: false

                }],
                dataBound: function (e) {
                    var columns = e.sender.columns;
                    var columnIndexDebet = this.wrapper.find(".k-grid-header [data-field=" + "Debet" + "]").index();
                    var columnIndexCredit = this.wrapper.find(".k-grid-header [data-field=" + "Credit" + "]").index();
                    var columnIndexSaldo = this.wrapper.find(".k-grid-header [data-field=" + "Saldo" + "]").index();
                    dataView = this.dataSource.view();
                    for (var i = 0; i < dataView.length; i++) {

                        if (dataView[i].Main_Acc_Code == 0) {
                            var uid = dataView[i].uid;
                            $("#gridMasterCOA").find("tr[data-uid=" + uid + "]").addClass("greyBackground");
                        }
                        if (dataView[i].Debet < 0) {
                            var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                            var cell1 = row1.children().eq(columnIndexDebet);
                            cell1.addClass('redValue')
                        }
                        if (dataView[i].Credit < 0) {
                            var row2 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                            var cell2 = row2.children().eq(columnIndexCredit);
                            cell2.addClass('redValue')
                        }
                        if (dataView[i].Saldo < 0) {
                            var row2 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                            var cell2 = row2.children().eq(columnIndexSaldo);
                            cell2.addClass('redValue')
                        }
                    }
                }
            })
        })

    },
    showDetail: function (ID, accname, acccode, parent) {
        mCoa.showDetailCoa(true)
        mCoa.HiddenFilterCoa(true)
        mCoa.detailCoaCode(acccode)
        mCoa.detailCoaName(accname)
        mCoa.detailCoaParent(parent)
        mCoa.textSearch("")
        $("#textSearch").val("")
        coa.renderGridDetailCoa(function(){
            coa.gridDetailCoa()
        })
    },
    BacktoCoa: function () {
        mCoa.HiddenFilterCoa(false)
        mCoa.showDetailCoa(false)
        mCoa.textSearch("")
        $("#textSearch").val("")
    },

    renderGridDetailCoa: function (callback) {
        model.Processing(true)
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
        var url = "/master/getdatadetailcoa"
        var param = {}
        if (mCoa.DateFilter() == true) {
            param = {
                DateStart: dateStart,
                DateEnd: dateEnd,
                Filter: true,
                Accountcode: mCoa.detailCoaCode(),
                ParentCode: mCoa.detailCoaParent()
            }
        } else {
            param = {
                Filter: false,
                Accountcode: mCoa.detailCoaCode(),
                ParentCode: mCoa.detailCoaParent()
            }
        }
        ajaxPost(url, param, function (res) {
            model.Processing(false)
            var Begining = res.Begining
            var Data = res.Data

            Data = _.sortBy(Data, [function (o) {
                return o.PostingDate
            }])
            var SaldoCalculate = 0.0
            for (i in Data) {
                if (i == 0) {
                    var begin = Begining + Data[i].Debet
                    Data[i].Saldo = begin - Data[i].Credit
                    SaldoCalculate = Data[i].Saldo
                }
                if (i > 0) {
                    var begin = SaldoCalculate + Data[i].Debet
                    Data[i].Saldo = begin - Data[i].Credit
                    SaldoCalculate = Data[i].Saldo
                }
            }
            mCoa.dataMasterDetailCoa($.extend(true, [], Data))
            mCoa.dataDetailCoa($.extend(true, [], Data))
            callback();
        })
    },
    gridDetailCoa : function(){
        var data = mCoa.dataDetailCoa()
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
            width: 90
        }, {
            field: 'DocumentNumber',
            title: 'Document Number',
            template: function (d) {
                // console.log(d)
                return '<a class="onclickCOA" onclick="coa.linkToJournal(\'' + d.IdJournal + '\')">' + d.DocumentNumber + '</a>'
            },
            width: 80
        }, {
            field: 'Department',
            title: 'Department',
            width: 70
        },{
            field: 'Description',
            title: 'Description',
            width: 90,
            footerTemplate: "<div style='text-align:right;font-size:15px'>Total :</div>"
        }, {
            field: 'Debet',
            title: 'Debit',
            width: 90,
            template: "#=ChangeToRupiah(Debet)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        }, {
            field: 'Credit',
            title: 'Credit',
            width: 90,
            template: "#=ChangeToRupiah(Credit)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        }, {
            field: 'Saldo',
            title: 'Saldo',
            width: 90,
            template: "#=ChangeToRupiah(Saldo)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= coa.saldoCalculation() #</div>",
            attributes: {
                style: "text-align:right;"
            }
        }, {
            field: 'User',
            title: 'User',
            width: 80
        }]
        $('#gridDetailCoa').kendoGrid({
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
                pageSize: 50,
            },
            pageable: {
                pageSizes: [50, 100, 500, "all"],           
                refresh: true,
                // pageSizes: true,
                buttonCount:5 
            },
            height: 500,
            filterable: {
                extra : false,
                operator : {
                    string : {
                        contains: "Contains",
                    }
                }
            },
            scrollable: true,
            columns: columns,
            dataBound: function (e) {
                var columns = e.sender.columns;
                var columnIndexDebet = this.wrapper.find(".k-grid-header [data-field=" + "Debet" + "]").index();
                var columnIndexCredit = this.wrapper.find(".k-grid-header [data-field=" + "Credit" + "]").index();
                var columnIndexSaldo = this.wrapper.find(".k-grid-header [data-field=" + "Saldo" + "]").index();
                dataView = this.dataSource.view();
                for (var i = 0; i < dataView.length; i++) {
                    if (dataView[i].Debet < 0) {
                        var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                        var cell1 = row1.children().eq(columnIndexDebet);
                        cell1.addClass('redValue')
                    }
                    if (dataView[i].Credit < 0) {
                        var row2 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                        var cell2 = row2.children().eq(columnIndexCredit);
                        cell2.addClass('redValue')
                    }
                    if (dataView[i].Saldo < 0) {
                        var row2 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                        var cell2 = row2.children().eq(columnIndexSaldo);
                        cell2.addClass('redValue')
                    }
                }
            }
        })
    },
    saldoCalculation: function () {
        var data = $('#gridDetailCoa').data('kendoGrid').dataSource.options.data
        var sumDebit = _.sumBy(data, 'Debet')
        var sumCredt = _.sumBy(data, 'Credit')
        var Saldo = sumDebit - sumCredt
        return kendo.toString(Saldo, 'n')
    },
    FilterAccCode: function (element) {
        element.kendoDropDownList({
            dataSource: mCoa.DataACC_Code(),
            filter: "startswith",
            optionLabel: "--Select Acc.Code--"
        });
    },
    RefreshDataByDate: function () {
        mCoa.DateFilter(true)
        coa.GetDataCoa()
    },
    ClearModalAdd: function () {
        $("#acccode").val("")
        $("#accname").val("")
        $("#mainacccode").val("")
        mCoa.DebetCredit("")
        mCoa.categoryCOA("")
    },
    SaveDataCOA: function () {
        var acccode = $("#acccode").val()
        if (acccode != parseInt(acccode, 10)) {
            return swal({
                title:"Error!",
                text: "Account code is Empty!",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
        var mainaccount = $("#mainacccode").val()
        if (mainaccount != parseInt(mainaccount, 10)) {
            return swal({
                title:"Error!",
                text: "Main account code is Empty!",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
        if (isNumeric(acccode) == false || isNumeric(mainaccount) == false) {
            return swal({
                title:"Error!",
                text: "Account Code or Main Code is Empty!",
                type: "error",
                confirmButtonColor: "#3da09a"
            });
        }
        var accname = $("#accname").val().trim().toUpperCase()
        var debetcredit = mCoa.DebetCredit()
        var category = mCoa.categoryCOA()
        if (accname == "" || acccode == "" || debetcredit == undefined || category == undefined || debetcredit == "" || category == "") {
            swal({
                title:"Error!",
                text:  "All field must filled",
                type: "error",
                confirmButtonColor: "#3da09a"});
        } else {
            var param = {
                AccCode: parseInt(acccode),
                MainAccCode: parseInt(mainaccount),
                AccName: accname.toUpperCase(),
                DebetOrCredit: debetcredit,
                Category: category
            }
            var url = "/master/insertnewcoa"
            model.Processing(true)
            ajaxPost(url, param, function (e) {
                model.Processing(false)
                if (e.status == true) {
                    swal({
                        title: "Success!", 
                        text: e.Message,
                        type: "success",
                        confirmButtonColor: '#3da09a',
                        confirmButtonText: "OK"});
                    $('#AddNewModal').modal('hide');
                    mCoa.DateFilter(false)
                    coa.GetDataCoa()
                    coa.ClearModalAdd()
                } else {
                    swal("Error!", e.Message, "error");
                }
            })

        }
    },
    ImportExcel: function () {
        event.stopPropagation();
        event.preventDefault();

        model.Processing(true);

        if ($('input[type=file]')[0].files[0] == undefined) {
            swal('Error', 'Please select a file to upload!', 'error');
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
                url: '/master/uploadfiles',
                data: data,
                cache: false,
                contentType: false,
                processData: false,
                type: 'POST',
                success: function (data) {
                    swal('Success', 'Data has been uploaded successfully!', 'success');
                    model.Processing(false);
                    $('#fDok').val('');
                    $('#fInfo').val('');
                    $('#ImportModal').modal('hide');
                    coa.GetDataCoa()
                }
            });
        } else {
            swal('Error', 'Please select a file to upload!', 'error');
            model.Processing(false);
        }
    }
}
coa.linkToJournal = function(idjournal){
    window.location.assign("/transaction/journal?id="+idjournal)
}
function isNumeric(n) {
    return !isNaN(parseFloat(n)) && isFinite(n);
}

mCoa.search = function() {
    if (mCoa.showDetailCoa()== false){
        mCoa.DateFilter(true)
        coa.GetDataCoa(function () {
            coa.renderGridDetailCoa()
        })
    }else{
        var data = mCoa.dataMasterDetailCoa()
        var result = _.filter(data, function(v){
            return v.Description.indexOf(mCoa.textSearch()) >= 0 ||v.DocumentNumber.indexOf(mCoa.textSearch()) >= 0||v.Acc_Name.indexOf(mCoa.textSearch()) >= 0
            ||v.User.indexOf(mCoa.textSearch()) >= 0
        })
        $("#MasterGridRole").html("");
        mCoa.dataDetailCoa(result)
        coa.gridDetailCoa()
    }
}

function ValidateFiles(elm) {
    var file = elm.val();
    var ext = file.split(".");
    ext = ext[ext.length - 1].toLowerCase();
    var arrayExtensions = ["xls", "xlsx", "csv"];

    if (arrayExtensions.lastIndexOf(ext) == -1) {
        // alert("Invalid file extension. Please upload an excel file format!");
        swal('Error', 'Invalid file extension. Please upload an excel file format!', 'error');
        elm.val('');
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

$(document).ready(function () {
    $(document).on('change', ':file', function () {
        ValidateFiles($(this));
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
    coa.DateShowing()
    coa.GetDataCoa()

    $("#acccode").keypress(function (event) {
        var ew = event.which;
        if (48 <= ew && ew <= 57) return true;
        if (65 <= ew && ew <= 90) return true;
        if (97 <= ew && ew <= 122) return true;
        return false;
    });
    $("#mainacccode").keypress(function (event) {
        var ew = event.which;
        if (48 <= ew && ew <= 57) return true;
        if (65 <= ew && ew <= 90) return true;
        if (97 <= ew && ew <= 122) return true;
        return false;
    });
    $("#accname").keypress(function (event) {
        var ew = event.which;
        if (ew == 32) return true;
        if (48 <= ew && ew <= 57) return true;
        if (65 <= ew && ew <= 90) return true;
        if (97 <= ew && ew <= 122) return true;
        return false;
    });

});
mCoa.ExportToPdf = function(){
    var Data = $("#gridMasterCOA").data("kendoGrid").dataSource.options.data
    for(i in Data){
        Data[i].ACC_Code = parseInt(Data[i].ACC_Code)
    }
    Data = _.sortBy(Data, function(o){return o.ACC_Code})
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var url = "/master/exportpdfcoa"
    var param = {
        Filter : mCoa.DateFilter(),
        DateStart: dateStart,
        DateEnd: dateEnd,
        Data : Data,
    }
    model.Processing(true)
    ajaxPost(url, param, function(e){
        window.open('/res/docs/master/' + e.Data, '_blank');
        model.Processing(false)
    })
}

mCoa.exportExcel = function () {
    $("#gridMasterCOA").getKendoGrid().saveAsExcel();
}
mCoa.onChangeDateStart = function(val){
    if (val.getTime()>mCoa.dateEnd().getTime()){
        mCoa.dateEnd(val)
    }
}
mCoa.init = function () {
    mCoa.getDateNow()
}
$(function () {
    mCoa.init()

})