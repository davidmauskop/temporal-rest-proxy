var R=Object.defineProperty;var E=Object.getOwnPropertySymbols;var q=Object.prototype.hasOwnProperty,A=Object.prototype.propertyIsEnumerable;var h=(n,o,c)=>o in n?R(n,o,{enumerable:!0,configurable:!0,writable:!0,value:c}):n[o]=c,w=(n,o)=>{for(var c in o||(o={}))q.call(o,c)&&h(n,c,o[c]);if(E)for(var c of E(o))A.call(o,c)&&h(n,c,o[c]);return n};import{n as P}from"./notifications-939cb4a5.js";const S=n=>{console.error(n),typeof n=="string"&&P.add("error",n),n instanceof Error&&P.add("error",n.message)},L=n=>typeof n=="function",U=(n,o)=>{const c=o instanceof URLSearchParams;return o&&!c&&(o=new URLSearchParams(o)),o?`${n}?${o}`:n},M=n=>n.message!==void 0&&typeof n.message=="string",V=async(n,o={},c=10)=>{const{params:u={},request:i=fetch,token:f,shouldRetry:x=!1,onError:d,retryInterval:y=5e3}=o;let{options:$}=o;const T=f?{next_page_token:f}:{},_=new URLSearchParams(w(w({},u),T)),g=U(n,_);try{$=v($);const s=await i(g,$),l=await s.json(),{status:k,statusText:I}=s;if(!s.ok)if(d&&L(d))d({status:k,statusText:I,body:l});else throw new Error(`${k}: ${I}`);return l}catch(s){if(S(s),x&&c>0)return new Promise(l=>{setTimeout(()=>{l(V(n,o,c-1))},y)})}},v=n=>{const o=w({credentials:"include"},n);return o.headers=O(n==null?void 0:n.headers),o},O=n=>{n||(n={});const o="_csrf=",c="X-CSRF-TOKEN";try{let i=document.cookie.split(";").find(f=>f.includes(o));i&&!n[c]&&(i=i.trim().slice(o.length),n[c]=i)}catch(u){console.error(u)}return n};let t="";t.endsWith("/")&&(t=t.slice(0,-1));const D=n=>n.split("/").map(o=>encodeURIComponent(o)).join("/"),F=n=>(n.startsWith("/")&&(n=n.slice(1)),`${t}/api/v1/${D(n)}`);function j(n,o){const c={cluster:"/cluster",settings:"/settings",user:"/me",namespaces:"/namespaces","task-queue":`/namespaces/${o==null?void 0:o.namespace}/task-queues/${o==null?void 0:o.queue}`,workflows:`/namespaces/${o==null?void 0:o.namespace}/workflows`,"workflows.open":`/namespaces/${o==null?void 0:o.namespace}/workflows/open`,"workflows.closed":`/namespaces/${o==null?void 0:o.namespace}/workflows/closed`,workflow:`/namespaces/${o==null?void 0:o.namespace}/workflows/${o==null?void 0:o.executionId}/executions/${o==null?void 0:o.runId}`,"workflow.terminate":`/namespaces/${o==null?void 0:o.namespace}/workflows/${o==null?void 0:o.executionId}/executions/${o==null?void 0:o.runId}/terminate`,events:`/namespaces/${o==null?void 0:o.namespace}/workflows/${o==null?void 0:o.executionId}/executions/${o==null?void 0:o.runId}/events`,query:`/namespaces/${o==null?void 0:o.namespace}/workflows/${o==null?void 0:o.executionId}/executions/${o==null?void 0:o.runId}/query`};return F(c[n])}export{j as a,L as b,S as h,M as i,V as r};
