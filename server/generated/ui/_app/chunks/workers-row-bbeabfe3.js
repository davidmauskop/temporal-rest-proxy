import{r as ee,a as te}from"./route-for-api-b4766198.js";import{S as se,i as le,s as ae,I as _e,e as b,t as H,k as L,c as T,a as g,g as S,d as v,n as O,b as I,f as re,J as a,K as de,L as pe,M as me,x as A,u as Q,h as oe,r as ne,w as ce,N as R,aM as ie,j as q,m as Y,o as M,Q as N,v as B,aB as ue}from"./vendor-07bbbbe5.js";import{f as fe}from"./format-date-460c9f8a.js";async function Ae(r,e=fetch){const s=await ee(te("task-queue",r),{request:e,params:{taskQueueType:"1"}}),t=await ee(te("task-queue",r),{request:e,params:{taskQueueType:"2"}});t.pollers.forEach(i=>{i.taskQueueTypes=["ACTIVITY"]}),s.pollers.forEach(i=>{i.taskQueueTypes=["WORKFLOW"]});const l=i=>(_,k)=>{const $=_[k.identity]||{lastAccessTime:void 0,taskQueueTypes:[]};return _[k.identity]={lastAccessTime:!$.lastAccessTime||$.lastAccessTime<k.lastAccessTime?k.lastAccessTime:$.lastAccessTime,taskQueueTypes:$.taskQueueTypes.concat([i])},_};return t.pollers.filter(i=>s.pollers.some(_=>{if(i.identity===_.identity)return i.taskQueueTypes=[..._.taskQueueTypes,...i.taskQueueTypes],i})),t.pollers.reduce(l("ACTIVITY"),s.pollers.reduce(l("WORKFLOW"),{})),{pollers:t.pollers,taskQueueStatus:t.taskQueueStatus}}function ve(r){let e,s,t,l,i,_,k,$,w,D,W,h,x,m,d;const C=r[1].default,u=_e(C,r,r[0],null);return{c(){e=b("section"),s=b("div"),t=b("div"),l=H("ID"),i=L(),_=b("div"),k=H("Last Accessed"),$=L(),w=b("div"),D=H("Workflow Table Handler"),W=L(),h=b("div"),x=H("Activity Handler"),m=L(),u&&u.c(),this.h()},l(p){e=T(p,"SECTION",{class:!0});var n=g(e);s=T(n,"DIV",{class:!0});var f=g(s);t=T(f,"DIV",{class:!0});var E=g(t);l=S(E,"ID"),E.forEach(v),i=O(f),_=T(f,"DIV",{class:!0});var j=g(_);k=S(j,"Last Accessed"),j.forEach(v),$=O(f),w=T(f,"DIV",{class:!0});var y=g(w);D=S(y,"Workflow Table Handler"),y.forEach(v),W=O(f),h=T(f,"DIV",{class:!0});var K=g(h);x=S(K,"Activity Handler"),K.forEach(v),f.forEach(v),m=O(n),u&&u.l(n),n.forEach(v),this.h()},h(){I(t,"class","w-3/12 text-left"),I(_,"class","w-3/12 text-left"),I(w,"class","w-3/12 text-left"),I(h,"class","w-2/12 text-left"),I(s,"class","bg-gray-200 flex flex-row p-2"),I(e,"class","flex flex-col border-2 border-gray-300 w-full rounded-lg")},m(p,n){re(p,e,n),a(e,s),a(s,t),a(t,l),a(s,i),a(s,_),a(_,k),a(s,$),a(s,w),a(w,D),a(s,W),a(s,h),a(h,x),a(e,m),u&&u.m(e,null),d=!0},p(p,[n]){u&&u.p&&(!d||n&1)&&de(u,C,p,p[0],d?me(C,p[0],n,null):pe(p[0]),null)},i(p){d||(A(u,p),d=!0)},o(p){Q(u,p),d=!1},d(p){p&&v(e),u&&u.d(p)}}}function ke(r,e,s){let{$$slots:t={},$$scope:l}=e;return r.$$set=i=>{"$$scope"in i&&s(0,l=i.$$scope)},[l,t]}class Qe extends se{constructor(e){super();le(this,e,ke,ve,ae,{})}}function $e(r){let e,s;return e=new R({props:{icon:ie,color:"black"}}),{c(){q(e.$$.fragment)},l(t){Y(e.$$.fragment,t)},m(t,l){M(e,t,l),s=!0},p:N,i(t){s||(A(e.$$.fragment,t),s=!0)},o(t){Q(e.$$.fragment,t),s=!1},d(t){B(e,t)}}}function he(r){let e,s;return e=new R({props:{icon:ue,color:"blue"}}),{c(){q(e.$$.fragment)},l(t){Y(e.$$.fragment,t)},m(t,l){M(e,t,l),s=!0},p:N,i(t){s||(A(e.$$.fragment,t),s=!0)},o(t){Q(e.$$.fragment,t),s=!1},d(t){B(e,t)}}}function be(r){let e,s;return e=new R({props:{icon:ie,color:"black"}}),{c(){q(e.$$.fragment)},l(t){Y(e.$$.fragment,t)},m(t,l){M(e,t,l),s=!0},p:N,i(t){s||(A(e.$$.fragment,t),s=!0)},o(t){Q(e.$$.fragment,t),s=!1},d(t){B(e,t)}}}function Te(r){let e,s;return e=new R({props:{icon:ue,color:"blue"}}),{c(){q(e.$$.fragment)},l(t){Y(e.$$.fragment,t)},m(t,l){M(e,t,l),s=!0},p:N,i(t){s||(A(e.$$.fragment,t),s=!0)},o(t){Q(e.$$.fragment,t),s=!1},d(t){B(e,t)}}}function ge(r){let e,s,t=r[0].identity+"",l,i,_,k,$,w=fe(r[0].lastAccessTime)+"",D,W,h,x,m,d,C,u,p,n,f,E;const j=[he,$e],y=[];function K(c,o){return(x==null||o&1)&&(x=!!c[0].taskQueueTypes.includes("WORKFLOW")),x?0:1}m=K(r,-1),d=y[m]=j[m](r);const J=[Te,be],V=[];function z(c,o){return(p==null||o&1)&&(p=!!c[0].taskQueueTypes.includes("ACTIVITY")),p?0:1}return n=z(r,-1),f=V[n]=J[n](r),{c(){e=b("article"),s=b("div"),l=H(t),i=L(),_=b("div"),k=b("h3"),$=b("p"),D=H(w),W=L(),h=b("div"),d.c(),C=L(),u=b("div"),f.c(),this.h()},l(c){e=T(c,"ARTICLE",{class:!0});var o=g(e);s=T(o,"DIV",{class:!0});var F=g(s);l=S(F,t),F.forEach(v),i=O(o),_=T(o,"DIV",{class:!0});var P=g(_);k=T(P,"H3",{});var G=g(k);$=T(G,"P",{});var U=g($);D=S(U,w),U.forEach(v),G.forEach(v),P.forEach(v),W=O(o),h=T(o,"DIV",{class:!0});var X=g(h);d.l(X),X.forEach(v),C=O(o),u=T(o,"DIV",{class:!0});var Z=g(u);f.l(Z),Z.forEach(v),o.forEach(v),this.h()},h(){I(s,"class","links w-3/12 text-left"),I(_,"class","links w-3/12 text-left"),I(h,"class","w-3/12 text-left"),I(u,"class","w-3/12 text-left"),I(e,"class","w-full h-full flex flex-row border-b-2 last:border-b-0 no-underline p-2")},m(c,o){re(c,e,o),a(e,s),a(s,l),a(e,i),a(e,_),a(_,k),a(k,$),a($,D),a(e,W),a(e,h),y[m].m(h,null),a(e,C),a(e,u),V[n].m(u,null),E=!0},p(c,[o]){(!E||o&1)&&t!==(t=c[0].identity+"")&&oe(l,t),(!E||o&1)&&w!==(w=fe(c[0].lastAccessTime)+"")&&oe(D,w);let F=m;m=K(c,o),m===F?y[m].p(c,o):(ne(),Q(y[F],1,1,()=>{y[F]=null}),ce(),d=y[m],d?d.p(c,o):(d=y[m]=j[m](c),d.c()),A(d,1),d.m(h,null));let P=n;n=z(c,o),n===P?V[n].p(c,o):(ne(),Q(V[P],1,1,()=>{V[P]=null}),ce(),f=V[n],f?f.p(c,o):(f=V[n]=J[n](c),f.c()),A(f,1),f.m(u,null))},i(c){E||(A(d),A(f),E=!0)},o(c){Q(d),Q(f),E=!1},d(c){c&&v(e),y[m].d(),V[n].d()}}}function we(r,e,s){let{poller:t}=e;return r.$$set=l=>{"poller"in l&&s(0,t=l.poller)},[t]}class De extends se{constructor(e){super();le(this,e,we,ge,ae,{poller:0})}}export{Qe as W,De as a,Ae as g};