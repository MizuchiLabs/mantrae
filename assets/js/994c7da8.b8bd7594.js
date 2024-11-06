"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[764],{3988:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>l,contentTitle:()=>i,default:()=>d,frontMatter:()=>o,metadata:()=>a,toc:()=>c});const a=JSON.parse('{"id":"usage/agents","title":"Agents","description":"Agents extend the functionality of your Traefik setup by allowing you to collect Docker container information from machines that don\u2019t have Traefik installed. This way, you can apply Traefik labels on containers across multiple hosts, and the agent will sync this data with your main Traefik instance.","source":"@site/docs/usage/03-agents.md","sourceDirName":"usage","slug":"/usage/agents","permalink":"/mantrae/docs/usage/agents","draft":false,"unlisted":false,"tags":[],"version":"current","sidebarPosition":3,"frontMatter":{"sidebar_position":3},"sidebar":"tutorialSidebar","previous":{"title":"DNS","permalink":"/mantrae/docs/usage/dns"},"next":{"title":"Environment","permalink":"/mantrae/docs/usage/environment"}}');var s=t(6070),r=t(5658);const o={sidebar_position:3},i="Agents",l={},c=[{value:"How Agents Work",id:"how-agents-work",level:2},{value:"Setting Up an Agent",id:"setting-up-an-agent",level:2},{value:"Step 1: Set the Server Address",id:"step-1-set-the-server-address",level:3},{value:"Step 2: Generate and Copy the Agent Token",id:"step-2-generate-and-copy-the-agent-token",level:3},{value:"Step 3: Run the Mantrae Agent",id:"step-3-run-the-mantrae-agent",level:3},{value:"Use Case: Traefik Labels on Remote Hosts",id:"use-case-traefik-labels-on-remote-hosts",level:2}];function h(e){const n={blockquote:"blockquote",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",hr:"hr",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,r.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.header,{children:(0,s.jsx)(n.h1,{id:"agents",children:"Agents"})}),"\n",(0,s.jsx)(n.p,{children:"Agents extend the functionality of your Traefik setup by allowing you to collect Docker container information from machines that don\u2019t have Traefik installed. This way, you can apply Traefik labels on containers across multiple hosts, and the agent will sync this data with your main Traefik instance."}),"\n",(0,s.jsx)(n.h2,{id:"how-agents-work",children:"How Agents Work"}),"\n",(0,s.jsx)(n.p,{children:"An agent is a standalone binary that runs on any machine where you want to collect container information. Each agent:"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Collects Docker container metadata, including Traefik labels."}),"\n",(0,s.jsx)(n.li,{children:"Communicates with the Mantrae server, sending back container info for unified management."}),"\n",(0,s.jsx)(n.li,{children:"Regularly renews its access token to ensure a secure, persistent connection."}),"\n"]}),"\n",(0,s.jsx)(n.h2,{id:"setting-up-an-agent",children:"Setting Up an Agent"}),"\n",(0,s.jsx)(n.h3,{id:"step-1-set-the-server-address",children:"Step 1: Set the Server Address"}),"\n",(0,s.jsxs)(n.p,{children:["In the settings, specify the ",(0,s.jsx)(n.strong,{children:"server address"})," for your Mantrae server. This address must be accessible by the agent to ensure successful communication."]}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Example"}),": If Mantrae is hosted on a public IP or domain (e.g., ",(0,s.jsx)(n.code,{children:"https://mantrae.example.com"}),"), configure this as the server address so agents can connect reliably."]}),"\n"]}),"\n",(0,s.jsxs)(n.blockquote,{children:["\n",(0,s.jsxs)(n.p,{children:[(0,s.jsx)(n.strong,{children:"Note"}),": The agent will automatically renew its token at regular intervals, so you don\u2019t need to worry about re-authenticating it manually."]}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"step-2-generate-and-copy-the-agent-token",children:"Step 2: Generate and Copy the Agent Token"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["In the Mantrae UI, navigate to the ",(0,s.jsx)(n.strong,{children:"Agents"})," tab."]}),"\n",(0,s.jsx)(n.li,{children:"Locate the generated token for the agent which authorizes the agent to communicate with the Mantrae server."}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"step-3-run-the-mantrae-agent",children:"Step 3: Run the Mantrae Agent"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsx)(n.li,{children:"Download the agent binary for your platform."}),"\n",(0,s.jsxs)(n.li,{children:["Run the agent with the token using the following command:","\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"./mantrae-agent -token <your_agent_token>\n"})}),"\n"]}),"\n",(0,s.jsx)(n.li,{children:"Ensure that the machine running the agent has Docker installed, as it will gather container details from the Docker daemon."}),"\n"]}),"\n",(0,s.jsx)(n.h2,{id:"use-case-traefik-labels-on-remote-hosts",children:"Use Case: Traefik Labels on Remote Hosts"}),"\n",(0,s.jsx)(n.p,{children:"Once the agent is running, you can set Traefik labels on containers located on the agent\u2019s host machine as you normally would for Traefik. The agent collects these labels and sends them to the Mantrae server, where they are applied to the main Traefik instance."}),"\n",(0,s.jsx)(n.p,{children:"This setup allows you to centralize routing configurations across multiple hosts without installing Traefik on each machine."}),"\n",(0,s.jsx)(n.hr,{}),"\n",(0,s.jsx)(n.p,{children:"Using agents, you can scale your Traefik configuration effortlessly, managing containers across multiple machines from a single, centralized Mantrae server."})]})}function d(e={}){const{wrapper:n}={...(0,r.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(h,{...e})}):h(e)}},5658:(e,n,t)=>{t.d(n,{R:()=>o,x:()=>i});var a=t(758);const s={},r=a.createContext(s);function o(e){const n=a.useContext(r);return a.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function i(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:o(e.components),a.createElement(r.Provider,{value:n},e.children)}}}]);