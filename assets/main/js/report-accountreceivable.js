var ar = {}
ar.DatePageBar= ko.observable("")
ar.dataCustomer= ko.observableArray([])
ar.customerCode= ko.observable("")
ar.dataTypeReport= ko.observableArray([])
ar.typeReportValue = ko.observable("")
ar.dataSummary = ko.observableArray([])
ar.dataDetail = ko.observableArray([])
ar.dataDetailGrouping = ko.observableArray([])
ar.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    ar.DatePageBar(page)
}
ar.getDataTypeReport = function(){
    var data = [
        {text:"Summary", value:"Summary"},
        {text:"Detail", value:"Detail"}
    ]
    ar.dataTypeReport(data)
}
ar.getDataCustomer = function () {
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
        ar.dataCustomer(res.Data)
        model.Processing(false)
    })
}
ar.getSummary = function(callback){
    var url = "/report/getdataar"
    var date = $("#dateAr").data("kendoDatePicker").value()
    var param ={
        DateAr : moment(date).format("DD-MMM-YYYY"),
        Customer : ar.customerCode(),
        Type : ar.typeReportValue() 
    }
    ajaxPost(url, param, function(res){
        var data = res.Data
        var group = _.groupBy(data, function(res){return res.CustomerCode})
        data = _.sortBy(group, [function(o) { return o.Due_Date; }]);
        var newData = []
        for(i in data){
            var each = _.sortBy(data[i], [function(o) { return o.Due_Date; }]);
            var sum = {
                "CustomerCode":each[0].CustomerCode,
                "CustomerName":each[0].CustomerName,
                "DocNum":each[0].DocNum,
                "Due_Date":each[0].Due_Date,
                "InvoiceDate":each[0].InvoiceDate,
                "Term":each[0].Term,
                "Amount": _.sumBy(each, function(o) { return o.Amount; }),
                "Total_AR": _.sumBy(each, function(o) { return o.Total_AR; }),
                "Age1":_.sumBy(each, function(o) { return o.Age1; }),
                "Age2":_.sumBy(each, function(o) { return o.Age2; }),
                "Age3":_.sumBy(each, function(o) { return o.Age3; }),
                "Age4":_.sumBy(each, function(o) { return o.Age4; })
            }
            newData.push(sum)
        }
        ar.dataSummary(newData)
        callback();
    })
}
ar.renderGridSummary = function(){
    $("#gridAr").html("")
    var columns = [
        {
            title : "Customer",
            field : "CustomerName",
            width: 250            
        },{
            title : "Term",
            field : "Term",
            width: 70
        },{
            title : "Due Date",
            field : "Due_Date",
            template: function (e) {
                return moment(e.Due_Date).format("DD MMM YYYY")
            },
            footerTemplate: "<div style='text-align:right;font-size:15px'>Total :</div>"
        },{
            title : "Amount",
            field : "Amount",
            template: "#=ChangeToRupiah(Amount)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Total AR",
            field : "Total_AR",
            template: "#=ChangeToRupiah(Total_AR)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged 1 - 30",
            field : "Age1",
            template: "#=ChangeToRupiah(Age1)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged 31 - 60",
            field : "Age2",
            template: "#=ChangeToRupiah(Age2)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged 61 - 90",
            field : "Age3",
            template: "#=ChangeToRupiah(Age3)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged > 90",
            field : "Age4",
            template: "#=ChangeToRupiah(Age4)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        }
    ]
    $('#gridAr').kendoGrid({
        dataSource: {
            data: ar.dataSummary(),
            aggregate: [{
                field: "Amount",
                aggregate: "sum"
            }, {
                field: "Total_AR",
                aggregate: "sum"
            }, {
                field: "Age1",
                aggregate: "sum"
            },{
                field: "Age2",
                aggregate: "sum"
            },{
                field: "Age3",
                aggregate: "sum"
            },{
                field: "Age4",
                aggregate: "sum"
            }],
            sort: {
                field: 'CustomerName',
                dir: 'asc',
            },
            pageSize: 10,
        },
        pageable: {
            pageSizes: [10, 20, 50, "all"],           
            refresh: true,
            // pageSizes: true,
            buttonCount:5 
        },
        excel: {
            fileName: "report-araging-summary.xlsx"
        },
        pdf: {
            allPages: true,
            paperSize: "A4",
            margin: { top: "3cm", right: "1cm", bottom: "1cm", left: "1cm" },
            landscape: true,
            template: $("#page-template").html()
        },
        sortable: true,
        scrollable: true,
        columns: columns,
        excelExport: function (e) {
            var rows = e.workbook.sheets[0].rows;
            for (var ri = 0; ri < rows.length; ri++) {
                var row = rows[ri];
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "group-footer" || row.type == "footer") {
                        var text = $(cell.value).text();
                        cell.hAlign = "right";
                        cell.value = text
                    }
                    if (ci == 2){
                        if (row.type == "data") {
                            cell.value = moment(cell.value).format("DD MMM YYYY")
                        }
                    }
                    if (ci>= 3) {
                        if (row.type == "data") {
                            cell.value = parseFloat(cell.value)
                            cell.format = "#,##0.00_);(#,##0.00);0.00;"
                            // Set the alignment
                            cell.hAlign = "right";
                        }
                    }
                }
            }
        }
    })
}
ar.getDetail = function(callback){
    var url = "/report/getdataar"
    var date = $("#dateAr").data("kendoDatePicker").value()
    var param ={
        DateAr : moment(date).format("DD-MMM-YYYY"),
        Customer : ar.customerCode(),
        Type : ar.typeReportValue() 
    }
    ajaxPost(url, param, function(res){
        var data = res.Data
        ar.dataDetail(data)
        var newdata = []
        for(i in data){
            if (newdata.length==0){
                newdata.push({
                    CustomerCode : data[i].CustomerCode,
                    Customer : data[i].CustomerName,
                    Item : [data[i]]
                })
            }else{
                var find = _.find(newdata, function(o) { return o.CustomerCode == data[i].CustomerCode });
                if(find!=undefined){
                    for (j in newdata){
                        if(data[i].CustomerCode ==newdata[j].CustomerCode ){
                            newdata[j].Item.push(data[i])                        
                        }
                    }
                }else{
                    newdata.push({
                        CustomerCode : data[i].CustomerCode,
                        Customer : data[i].CustomerName,
                        Item : [data[i]]
                    })
                }
            }
        }
        ar.dataDetailGrouping(newdata)
        callback();
    })
}
ar.renderGridDetail = function(){
    $("#gridAr").html("")
    var columns = [
        {
            title : "Customer",
            field : "CustomerName",
            width: 250            
        },{
            title : "No. Doc",
            field : "DocNum",         
        },{
            title : "Invoice Date",
            field : "InvoiceDate",
            template: function (e) {
                return moment(e.InvoiceDate).format("DD MMM YYYY")
            },
        },{
            title : "Due Date",
            field : "Due_Date",
            template: function (e) {
                return moment(e.Due_Date).format("DD MMM YYYY")
            },
            footerTemplate: "<div style='text-align:right;font-size:15px'>Total All:</div>",
            groupFooterTemplate : function(e){
                var title = ""
                if(e.Age1.group!= undefined){
                    var text = e.Age1.group.value
                    var newtext = text.split("-");
                    title = newtext[1] 
                }
                return "<div style='text-align:right;font-size:15px'>"+title+" Total :</div>"
            }
            // groupFooterTemplate: "<div style='text-align:right;font-size:15px'>Total :</div>"
        },{
            title : "Amount",
            field : "Amount",
            template: "#=ChangeToRupiah(Amount)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged 1 - 30",
            field : "Age1",
            template: "#=ChangeToRupiah(Age1)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            groupFooterTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged 31 - 60",
            field : "Age2",
            template: "#=ChangeToRupiah(Age2)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            groupFooterTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged 61 - 90",
            field : "Age3",
            template: "#=ChangeToRupiah(Age3)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            groupFooterTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        },{
            title : "Aged > 90",
            field : "Age4",
            template: "#=ChangeToRupiah(Age4)#",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            groupFooterTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            attributes: {
                style: "text-align:right;"
            }
        }
    ]
    $('#gridAr').kendoGrid({
        dataSource: {
            data: ar.dataDetail(),
            aggregate: [{
                field: "Amount",
                aggregate: "sum"
            },{
                field: "Age1",
                aggregate: "sum"
            },{
                field: "Age2",
                aggregate: "sum"
            },{
                field: "Age3",
                aggregate: "sum"
            },{
                field: "Age4",
                aggregate: "sum"
            }],
            group: { field: "CustomerName", aggregates: [ 
                { field: "Amount", aggregate: "sum" },
                { field: "Age1", aggregate: "sum" },
                { field: "Age2", aggregate: "sum" },
                { field: "Age3", aggregate: "sum" },
                { field: "Age4", aggregate: "sum" },
            ] },
            sort: {
                field: 'CustomerName',
                dir: 'asc',
            },
            pdf: {
                allPages: true,
                paperSize: "A4",
                margin: { top: "3cm", right: "1cm", bottom: "1cm", left: "1cm" },
                landscape: true,
                template: $("#page-template").html()
            },
            pageSize: 10,
        },
        pageable: {
            pageSizes: [10, 20, 50, "all"],           
            refresh: true,
            // pageSizes: true,
            buttonCount:5 
        },
        excel: {
            fileName: "report-araging-detail.xlsx"
        },
        sortable: true,
        scrollable: true,
        columns: columns,
        excelExport: function (e) {
            var rows = e.workbook.sheets[0].rows;
            for (var ri = 0; ri < rows.length; ri++) {
                var row = rows[ri];
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "group-footer" || row.type == "footer") {
                        var text = $(cell.value).text();
                        cell.hAlign = "right";
                        cell.value = text
                    }
                    if (ci == 3 ||ci == 4){
                        if (row.type == "data") {
                            cell.value = moment(cell.value).format("DD MMM YYYY")
                        }
                    }
                    if (ci> 4) {
                        if (row.type == "data") {
                            cell.value = parseFloat(cell.value)
                            cell.format = "#,##0.00_);(#,##0.00);0.00;"
                            // Set the alignment
                            cell.hAlign = "right";
                        }
                    }
                }
            }
        }
    })
}
ar.generate = function(){
    if(ar.typeReportValue()=="Summary"){
        ar.getSummary(function(){
            ar.renderGridSummary()    
        })
    }else if(ar.typeReportValue()=="Detail"){
        ar.getDetail(function(){
            ar.renderGridDetail()
            ar.removeRowGroup()            
        })
    }else {
        return swal({
            title: "Warning!", 
            text: "Please check type report!", 
            type: "warning",
            confirmButtonColor: "#3da09a"})
    }
}
ar.removeRowGroup  = function(){
    $(".k-group-col,.k-group-cell").remove();
    $(".k-grid .k-icon,.k-grid .k-i-collapse").remove();
    var spanCells = $(".k-grouping-row").children("td");
    spanCells.attr("colspan", spanCells.attr("colspan") - 1);
}
ar.downloadExcel = function(){
    $("#gridAr").getKendoGrid().saveAsExcel();
}
ar.ExportToPdf = function(){
    // $('.k-grid').css('font-size','8px');
    // $('.k-header').css('font-size','8px');
    // $('.k-footer-template td div').css('font-size','8px'); 
    // $('.k-grid td').css('line-height','2em');
    // $("#gridAr").getKendoGrid().saveAsPDF();
    // setTimeout(function(){
    //     $('.k-grid').css('font-size','12px');
    //     $('.k-header').css('font-size','12px');
    //     $('.k-footer-template td div').css('font-size','15px'); 
    //     $('.k-grid td').css('line-height','2em');
    // },100)
    model.Processing(true)
    var date = $("#dateAr").data("kendoDatePicker").value()
    var param = {
        DateAr : moment(date).format("DD-MMM-YYYY"),
        Type : ar.typeReportValue(),
        DataDetail : ar.dataDetailGrouping(),
        DataSummary: ar.dataSummary()
    }
    ajaxPost("/report/exportpdfar", param, function (e) {
        model.Processing(false)
        var taborWindow = window.open('/res/docs/report/pdf/' + e.Data, '_blank');
        taborWindow.focus();
    })
}
ar.init = function(){
   ar.getDateNow()
   ar.getDataCustomer()
   ar.getDataTypeReport()
   ar.typeReportValue("Summary")
   ar.generate()

}
$(document).ready(function () {
   ar.init()
})