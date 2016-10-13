"use strict";
var Bot = (function () {
    function Bot(args) {
        this.args = args;
    }
    Bot.prototype.execute = function (cb) {
        throw "Nothing Implemented";
    };
    return Bot;
}());
exports.Bot = Bot;
