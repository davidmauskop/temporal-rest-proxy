import{an as a,ao as s,ap as i,aq as c}from"./vendor-c7e306c1.js";const r="H:mm:ss:SS a O ccc, LLL do";function u(e,n="UTC"){if(!e)return"";try{t(e)&&(e=f(e));const o=a(e);return n==="local"?s(o,r):n==="relative"?i(o)+" ago":c(o,"UTC",r)}catch{return""}}function f(e){if(!t(e))throw new TypeError("provided value is not a timestamp");const n=new Date(null);return n.setTime(Number(e.seconds)*1e3+e.nanos/1e3),n}function t(e){return typeof e=="object"?e.seconds!==void 0&&e.nanos!==void 0:!1}export{u as f};
