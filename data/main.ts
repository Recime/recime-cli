/// <reference path="node.d.ts" />

import util = require('util')

export class Bot {

  private args: any;

  constructor(args:Object){
    this.args = args;
  }

  execute(cb:any){
    throw "Nothing Implemented";
  }
}
