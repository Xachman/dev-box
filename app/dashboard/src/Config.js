export default class Config {

    static get relative() { return false}
    static get devBoxHost(){  return "192.168.1.151"}
    static get devBoxPort(){ return 9080 }
    static get protocol() { return "http" }

    static dockerHostUrl() {
        if(this.relative) {
            return "/";
        }
        return this.protocol+"://"+this.devBoxHost+":"+this.devBoxPort+"/";
    }
    static host() {
        return this.devBoxHost
    }

    static port() {
        return this.devBoxPort
    }
}