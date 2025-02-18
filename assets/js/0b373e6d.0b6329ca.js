"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[537],{8770:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>l,contentTitle:()=>c,default:()=>o,frontMatter:()=>t,metadata:()=>i,toc:()=>h});const i=JSON.parse('{"id":"usage/environment","title":"Environment Variables","description":"Mantrae provides several command-line flags and environment variables to configure the application. This guide details each option and its purpose.","source":"@site/docs/usage/04-environment.md","sourceDirName":"usage","slug":"/usage/environment","permalink":"/mantrae/docs/usage/environment","draft":false,"unlisted":false,"tags":[],"version":"current","sidebarPosition":4,"frontMatter":{"sidebar_position":4},"sidebar":"tutorialSidebar","previous":{"title":"Agents","permalink":"/mantrae/docs/usage/agents"},"next":{"title":"Backups & Restoration","permalink":"/mantrae/docs/usage/backups"}}');var d=r(6070),s=r(385);const t={sidebar_position:4},c="Environment Variables",l={},h=[{value:"Command-Line Arguments",id:"command-line-arguments",level:2},{value:"Environment Variables",id:"environment-variables-1",level:2},{value:"Core Configuration",id:"core-configuration",level:3},{value:"Server Configuration",id:"server-configuration",level:3},{value:"Admin Configuration",id:"admin-configuration",level:3},{value:"Email Configuration",id:"email-configuration",level:3},{value:"Database Configuration",id:"database-configuration",level:3},{value:"Backup Configuration",id:"backup-configuration",level:3},{value:"Traefik Configuration",id:"traefik-configuration",level:3},{value:"Background Jobs Configuration",id:"background-jobs-configuration",level:3},{value:"Example Usage",id:"example-usage",level:3},{value:"Important Notes",id:"important-notes",level:3}];function a(e){const n={code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",p:"p",pre:"pre",strong:"strong",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,s.R)(),...e.components};return(0,d.jsxs)(d.Fragment,{children:[(0,d.jsx)(n.header,{children:(0,d.jsx)(n.h1,{id:"environment-variables",children:"Environment Variables"})}),"\n",(0,d.jsx)(n.p,{children:"Mantrae provides several command-line flags and environment variables to configure the application. This guide details each option and its purpose."}),"\n",(0,d.jsx)(n.h2,{id:"command-line-arguments",children:"Command-Line Arguments"}),"\n",(0,d.jsx)(n.p,{children:"You can use the following flags to customize the behavior of Mantrae:"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Flag"}),(0,d.jsx)(n.th,{children:"Type"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"-version"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"bool"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"false"})}),(0,d.jsx)(n.td,{children:"Prints the current version of Mantrae and exits."})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"-update"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"bool"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"false"})}),(0,d.jsx)(n.td,{children:"Updates Mantrae to the latest version. (Doesn't work inside a container)"})]})]})]}),"\n",(0,d.jsx)(n.h2,{id:"environment-variables-1",children:"Environment Variables"}),"\n",(0,d.jsx)(n.p,{children:"Environment variables can be used to set up Mantrae and configure its settings. Below is a list of the supported environment variables."}),"\n",(0,d.jsx)(n.h3,{id:"core-configuration",children:"Core Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsx)(n.tbody,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SECRET"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"Secret key required for secure access. Required!"})]})})]}),"\n",(0,d.jsx)(n.h3,{id:"server-configuration",children:"Server Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SERVER_HOST"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"0.0.0.0"})}),(0,d.jsx)(n.td,{children:"Host address the server will bind to"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SERVER_PORT"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"3000"})}),(0,d.jsx)(n.td,{children:"Port which Mantrae will listen on"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SERVER_URL"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"http://127.0.0.1"})}),(0,d.jsx)(n.td,{children:"The public URL of the Mantrae server for agent connections"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SERVER_BASIC_AUTH"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"false"})}),(0,d.jsx)(n.td,{children:"Enables basic authentication for the Mantrae server"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SERVER_ENABLE_AGENT"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"true"})}),(0,d.jsx)(n.td,{children:"Enables the Mantrae agent functionality"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"SERVER_LOG_LEVEL"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"info"})}),(0,d.jsx)(n.td,{children:"Log verbosity level"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"admin-configuration",children:"Admin Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"ADMIN_USERNAME"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"admin"})}),(0,d.jsx)(n.td,{children:"Admin user username"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"ADMIN_EMAIL"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"admin@mantrae"})}),(0,d.jsx)(n.td,{children:"Admin user email"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"ADMIN_PASSWORD"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"Admin user password"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"email-configuration",children:"Email Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"EMAIL_HOST"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"localhost"})}),(0,d.jsx)(n.td,{children:"SMTP server host"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"EMAIL_PORT"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"587"})}),(0,d.jsx)(n.td,{children:"SMTP server port"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"EMAIL_USERNAME"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"SMTP server username"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"EMAIL_PASSWORD"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"SMTP server password"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"EMAIL_FROM"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"mantrae@localhost"})}),(0,d.jsx)(n.td,{children:"Sender email address"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"database-configuration",children:"Database Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"DB_TYPE"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"sqlite"})}),(0,d.jsxs)(n.td,{children:["Database type. Supported options: only ",(0,d.jsx)(n.code,{children:"sqlite"})," for now"]})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"DB_NAME"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"mantrae"})}),(0,d.jsx)(n.td,{children:"Database/file name"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"backup-configuration",children:"Backup Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"BACKUP_ENABLED"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"true"})}),(0,d.jsx)(n.td,{children:"Enable automatic backups"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"BACKUP_PATH"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"backups"})}),(0,d.jsx)(n.td,{children:"Directory for storing backups"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"BACKUP_INTERVAL"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"24h"})}),(0,d.jsx)(n.td,{children:"Interval between backups"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"BACKUP_KEEP"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"3"})}),(0,d.jsx)(n.td,{children:"Number of backups to keep"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"traefik-configuration",children:"Traefik Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"TRAEFIK_PROFILE"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"default"})}),(0,d.jsx)(n.td,{children:"Traefik profile name"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"TRAEFIK_URL"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"Traefik API URL"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"TRAEFIK_USERNAME"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"Traefik authentication username"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"TRAEFIK_PASSWORD"})}),(0,d.jsx)(n.td,{}),(0,d.jsx)(n.td,{children:"Traefik authentication password"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"TRAEFIK_TLS"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"false"})}),(0,d.jsx)(n.td,{children:"Enable TLS for Traefik"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"background-jobs-configuration",children:"Background Jobs Configuration"}),"\n",(0,d.jsxs)(n.table,{children:[(0,d.jsx)(n.thead,{children:(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.th,{children:"Variable"}),(0,d.jsx)(n.th,{children:"Default"}),(0,d.jsx)(n.th,{children:"Description"})]})}),(0,d.jsxs)(n.tbody,{children:[(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"BACKGROUND_JOBS_TRAEFIK"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"20"})}),(0,d.jsx)(n.td,{children:"Traefik background job interval"})]}),(0,d.jsxs)(n.tr,{children:[(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"BACKGROUND_JOBS_DNS"})}),(0,d.jsx)(n.td,{children:(0,d.jsx)(n.code,{children:"300"})}),(0,d.jsx)(n.td,{children:"DNS background job interval"})]})]})]}),"\n",(0,d.jsx)(n.h3,{id:"example-usage",children:"Example Usage"}),"\n",(0,d.jsx)(n.p,{children:"To run Mantrae with custom environment variables:"}),"\n",(0,d.jsx)(n.pre,{children:(0,d.jsx)(n.code,{className:"language-bash",children:'export SECRET="your-secret-key"\nexport SERVER_PORT="4000"\nexport ADMIN_PASSWORD="secure-password"\n./mantrae\n'})}),"\n",(0,d.jsx)(n.h3,{id:"important-notes",children:"Important Notes"}),"\n",(0,d.jsxs)(n.ul,{children:["\n",(0,d.jsxs)(n.li,{children:[(0,d.jsx)(n.strong,{children:"SECRET"})," is a required environment variable and must be set; otherwise, the application will not start."]}),"\n",(0,d.jsxs)(n.li,{children:["Set ",(0,d.jsx)(n.strong,{children:"SERVER_URL"})," to the publicly accessible URL of Mantrae to ensure agents can connect to it."]}),"\n"]})]})}function o(e={}){const{wrapper:n}={...(0,s.R)(),...e.components};return n?(0,d.jsx)(n,{...e,children:(0,d.jsx)(a,{...e})}):a(e)}},385:(e,n,r)=>{r.d(n,{R:()=>t,x:()=>c});var i=r(758);const d={},s=i.createContext(d);function t(e){const n=i.useContext(s);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function c(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(d):e.components||d:t(e.components),i.createElement(s.Provider,{value:n},e.children)}}}]);