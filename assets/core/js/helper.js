var hasClass = function(el, className) {
    if (el.classList)
        return el.classList.contains(className)
    else
        return !!el.className.match(new RegExp('(\\s|^)' + className + '(\\s|$)'))
}

var addClass = function(el, className) {
    if (el.classList)
        el.classList.add(className)
    else if (!hasClass(el, className)) el.className += " " + className
}

var getDateFromString = function(str){
    var pattern = /(\d{1,4})([/\W+/g]|[/\W+/g](st|nd|rd|th)[/\W+/g]|(st|nd|rd|th)[/\W+/g])(\d{2}|.{3}|\b(?:January|February|March|April|May|June|July|August|September|October|November|December))([/\W+/g]|[/\W+/g]{2})(\d{4}|\d{2})/g
    var pattern2 = /(\d{3}|\b(?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May?|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?))[/\W+/g]*(\d{2})[/\W+/g]*(\d{1,4})/g
    str = str.replace('of ','')
    str = str.replace(/[^0-9a-z]/mgui, ' ')
    str = str.replace(/[^0-9a-z]/mgui, ' ')
    str = str.replace(/\n/g, " ")
    str = str.replace(/\s\s+/g, " ")

    //try pattern 1
    res = pattern.exec(str)
    var dateres,sorted,datedata,regexresult = null
    
    if (res != null){
        regexresult = _.reject(_.compact(res), function(v){ return /^\W+/g.test(v) == true })
        
        loop1:
        for(var i = 1; i < regexresult.length; i++){
            for(var j = 2; j < regexresult.length; j++){
                if( i == j) continue

                for(var k = 3; k < regexresult.length; k++){
                    if (((function() {
                        var o = /\b(st|nd|rd|th)/gi
                        return o.test(regexresult[i]) || o.test(regexresult[j]) || o.test(regexresult[k])
                    })())) continue

                    if (/\W/g.test(regexresult[i] + "" + regexresult[j] + "" + regexresult[k])) continue

                    if (i == k || j == k) continue

                    if( ( ! isNaN(parseInt(regexresult[i]))) && ( ! isNaN(parseInt(regexresult[j]))) && ( ! isNaN(parseInt(regexresult[k]))) ) {
                        sorted = _.sortBy( [parseInt(regexresult[i]), parseInt(regexresult[j]), parseInt(regexresult[k])] )
                        datedata = moment(sorted[0] + " " + sorted[1] + " " + sorted[2]).toDate()
                    } else {
                        datedata = moment(regexresult[i] + " " + regexresult[j] + " " + regexresult[k]).toDate()
                    }

                    if(datedata != "Invalid Date"){
                        dateres = {
                            regexres : regexresult,
                            date : datedata
                        }
                        break loop1
                    }
                }
            }
        }
    }else{
        //try pattern 2
        ress = pattern2.exec(str)
        console.log(ress)
        if (ress != null){
            regexresult = ress
            datedata = moment(ress[0]).toDate()
        }
    
        if(datedata != "Invalid Date"){
            dateres = {
                regexres : regexresult,
                date : datedata
            }
        }
    }
    
    console.log(dateres)
    return dateres
}

function getIrisan(strLong, strShort) {
    var doCheck = function (leftStr, rightStr) {
        for (var i = 1; i <= leftStr.length; i++) {
            var each = leftStr.slice(leftStr.length - i)
            if (rightStr.indexOf(each) == 0) {
                return each
            }
        }

        return ""
    }

    var leftResult = doCheck(strLong, strShort)
    if (leftResult != "") {
        return leftResult
    }

    var rightResult = doCheck(strShort, strLong)
    if (rightResult != "") {
        return rightResult
    }

    if (strLong.indexOf(strShort) !== -1) {
        return strShort
    }

    return ""
}

var readGetFromUrl = function(url = this.location.href) {
    var queryStart = url.indexOf("?") + 1,
        queryEnd   = url.indexOf("#") + 1 || url.length + 1,
        query = url.slice(queryStart, queryEnd - 1),
        pairs = query.replace(/\+/g, " ").split("&"),
        parms = {}, i, n, v, nv;

    if (query === url || query === "") return;

    for (i = 0; i < pairs.length; i++) {
        nv = pairs[i].split("=", 2);
        n = decodeURIComponent(nv[0]);
        v = decodeURIComponent(nv[1]);

        if (!parms.hasOwnProperty(n)) parms[n] = [];
        parms[n].push(nv.length === 2 ? v : null);
    }
    return parms;
}

var testRegex = function(str){
    var pattern = /\b(?:AFGHANISTAN|ANGOLA|AUSTRALIA|BAHRAIN|BANGLADESH|BOTSWANA|BRAZIL|BRUNEI DARUSSALAM|CAMEROON|CANADA|CHINA|COTE D'IVOIRE|FRANCE|GAMBIA|GERMANY|GHANA|GUERNSEY|HONG KONG|INDIA|INDONESIA|IRAN|IRAQ|IRELAND|ITALY|JAPAN|JORDAN|KENYA|KOREA|LEBANON|LUXEMBOURG|MALAYSIA|MAURITIUS|MACAU|NEPAL|NIGERIA|OMAN|PAKISTAN|PHILIPPINES|QATAR|SAUDI ARABIA|SIERRA LEONE|SINGAPORE|SOUTH AFRICA|SRI LANKA|SWEDEN|SWITZERLAND|TAIWAN|TANZANIA|THAILAND|TURKEY|UGANDA|UNITED ARAB EMIRATES|UNITED KINGDOM|UNITED STATES OF AMERICA|VIETNAM|ZAMBIA|ZIMBABWE)/

    return pattern.exec(str)

}

var toTitleCase = function(str){
    return str.replace(/\w\S*/g, function(txt){
        return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase()
    });
}

var getClusterTitleCase = function(text){
    return toTitleCase(text).replace(/\s/g, '')
}

var isInteger = function(value) {
  return /^\d+$/.test(value);
}

var getFuzzyResult = function(str){
    var result = ""
    var strword = str.toUpperCase()
    var masterdata = countryMaster.get("arr")
    for (var i = 0;i < masterdata.length;i++){
        if (fuzzysearch(masterdata[i],strword)){
            result = testRegex(strword)
        }
    }

    if (result == null){
        if (strword.indexOf("ENGLAND") != -1){
            result = "UNITED KINGDOM"
        }else if (strword.indexOf("IRELAND") != -1){
            result = "UNITED KINGDOM"
        }else if (strword.indexOf("WALES") != -1){
            result = "UNITED KINGDOM"
        }else if (strword.indexOf("SCOTLAND") != -1){
            result = "UNITED KINGDOM"
        }else{
            result = ""
        }
    }
    console.log(result)
    if( typeof result === 'object' ) {
        return result[0]
    }else{
        return result
    }
}