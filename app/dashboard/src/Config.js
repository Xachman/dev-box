export default class Config {
    static get devBoxHost(){  return "localhost"}
    static get devBoxPort(){ return 9080 }
    static get protocol() { return "http" }

    static dockerHostUrl() {
        return this.protocol+"://"+this.devBoxHost+":"+this.devBoxPort+"/";
    }
    static host() {
        return this.devBoxHost
    }

    static port() {
        return this.devBoxPort
    }
}