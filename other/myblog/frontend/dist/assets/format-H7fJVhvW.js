const a=new Intl.DateTimeFormat("zh-CN",{year:"numeric",month:"long",day:"numeric"});function n(e){const t=new Date(e);return Number.isNaN(t.getTime())?"日期未知":a.format(t)}export{n as f};
