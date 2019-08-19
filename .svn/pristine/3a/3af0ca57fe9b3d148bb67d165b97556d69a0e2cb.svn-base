function ChangeToRupiah(angka) {
    if (angka >= 0) {
        var TotString = kendo.toString(angka, "n");
        return TotString;
    } else {
        var TotminString = kendo.toString(Math.abs(angka), "n");
        return "(" + TotminString + ")";
    }
}

function FormatCurrency(angka) {
    var nilaiangka = kendo.toString(angka, "n");
    var nom = Number(nilaiangka.replace(/[^0-9\.]+/g, ""));
    return nom 
}

function MakeID() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 10; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}

function mongoObjectId() {
    var timestamp = (new Date().getTime() / 1000 | 0).toString(16);
    timestamp + 'xxxxxxxxxxxxxxxx'.replace(/[x]/g, function () {
        return (Math.random() * 16 | 0).toString(16);
    }).toLowerCase();

};

function validateNumber(event) {
    var key = window.event ? event.keyCode : event.which;
    if (event.keyCode === 8 || event.keyCode === 46) {
        return true;
    } else if ( key < 48 || key > 57 ) {
        return false;
    } else {
    	return true;
    }
};

