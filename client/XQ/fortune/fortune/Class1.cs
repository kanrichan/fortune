using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.IO;
using System.Threading.Tasks;
using System.Net;
using System.Security.Cryptography;
using XQ.Net.SDK.Attributes;
using XQ.Net.SDK.EventArgs;
using XQ.Net.SDK;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System.Runtime.InteropServices;

[assembly: XQPlugin("fortune-运势", "木理", "1.0.4", "项目地址 https://github.com/Yiwen-Chan/fortune")]

namespace fortune
{
    [Plugin]
    public class Fortune
    {
        public struct FromDataParm
        {
            public string Client;
            public string Version;
            public string Bot;
            public string Types;
            public string FromGroup;
            public string FromQQ;
            public string Ask;
            public string Limit;
        };

        public struct HeaderParm
        {
            public string Authkey;
            public string Au_time;
        }

        public static JObject OpenJson(string path)
        {
            StreamReader streamReader = new StreamReader(path);
            string text = streamReader.ReadToEnd();
            JObject json = (JObject)JsonConvert.DeserializeObject(text);
            return json;
        }

        public static void WriteJson(string path)
        {
            try
            {
                string json = @"{
    '默认': {
        '触发': '运势',
        '回复': '少女祈祷中......',
        '类型': '李清歌|碧蓝幻想|公主连结',
        '限制': '全局'
    },
    '单个群设置填群号': {
        '触发': '这里填关键词',
        '回复': '这里填收到关键词的回复',
        '类型': '这里是池子类型，多个池子用|分开',
        '限制': '每天限制一张，可填 全局 池子 关'
    },
    '00000000': {
        '触发': '抽签',
        '回复': '少女折寿中......',
        '类型': '车万',
        '限制': '池子'
    },
    '1048452984': {
        '触发': '运势测试',
        '回复': '收到命令！',
        '类型': '李清歌',
        '限制': '关'
    }
}";
                File.WriteAllText(path, json);
                XQAPI.OutPutLog("[fortune-运势] 检测到本插件为首次运行");
                XQAPI.OutPutLog("[fortune-运势] 配置已生成 " + path);
                XQAPI.OutPutLog("[fortune-运势] 若有需求可自行修改配置");
            }
            catch
            {
                XQAPI.OutPutLog("默认配置写入失败，请到项目地址提交issue");
            }
        }

        public static string ReadJson(JObject configJson, string section, string key, string def)
        {
            string result = "";
            try
            {
                result = configJson[section][key].ToString();
            }
            catch
            {
                result = def;
            }
            return result;
        }

        [EnableEvent]
        public static void enable(object sender, XQEventArgs e)
        {
            string configPath = Path.Combine(XQAPI.AppDir, "config.json");
            if (!File.Exists(@configPath))
            {
                WriteJson(@configPath);
            }
            XQAPI.OutPutLog("[fortune-运势] 项目地址 https://github.com/Yiwen-Chan/fortune");
            XQAPI.OutPutLog("[fortune-运势] 配置目录 " + configPath);
            XQAPI.OutPutLog("[fortune-运势] 特别感谢 fz6m https://github.com/fz6m/nonebot-plugin/tree/master/CQVortune");
            XQAPI.OutPutLog("[fortune-运势] 特别感谢 Lostdegree https://github.com/Lostdegree/Portune");
            XQAPI.OutPutLog("[fortune-运势] 想自定义运势背景并共享可加QQ群 1048452984 ");
        }

        [GroupMsgEvent]
        public static void onGroupMsg(object sender, XQAppGroupMsgEventArgs e)
        {
            string text = e.Message.Text;
            string robot = e.RobotQQ;
            string group = e.FromGroup.Id;
            string user = e.FromQQ.Id;

            string configPath = Path.Combine(XQAPI.AppDir, "config.json");
            JObject configJson = OpenJson(@configPath);

            if (ReadJson(configJson, group, "触发", "null") == "null")
            {
                group = "默认";
            }
            else if (ReadJson(configJson, group, "触发", "null") == "关")
            {
                return;
            }
            string trigger = ReadJson(configJson, group, "触发", "运势");
            string reply = ReadJson(configJson, group, "回复", "少女祈祷中......");
            string types = ReadJson(configJson, group, "类型", "李清歌");
            string limit = ReadJson(configJson, group, "限制", "池子");

            if (types.Contains("|"))
            {
                List<string> list = new List<string>(types.Split('|'));
                int length = list.ToArray().Length;
                Random ran = new Random();
                int RandKey = ran.Next(0, length);
                types = list[RandKey];
            }

            if (limit == "全局")
            {
                limit = "on";
            }
            else if (limit == "池子")
            {
                limit = "none";
            }
            else if (limit == "关")
            {
                limit = "off";
            }
            else
            {
                limit = "none";
            }

            if (text == trigger)
            {
                XQAPI.OutPutLog("[fortune-运势] 开始向服务器请求数据......");
                e.FromGroup.SendMessage(e.RobotQQ, reply);

                string aukey = "test";
                string apifortune = "127.0.0.1:8000";
                string apipic = "127.0.0.1:8000";

                FromDataParm fromDataParm;
                fromDataParm.Client = "xq";
                fromDataParm.Version = "4";
                fromDataParm.Bot = robot;
                fromDataParm.Types = types;
                fromDataParm.FromGroup = e.FromGroup.Id;
                fromDataParm.FromQQ = user;
                fromDataParm.Ask = text;
                fromDataParm.Limit = limit;

                string autime = get_autime();
                string authkey = get_authkey(aukey, autime);
                HeaderParm headerParm;
                headerParm.Authkey = authkey;
                headerParm.Au_time = autime;

                string fortuneJson = fortune(apifortune, fromDataParm, headerParm);

                JObject json = (JObject)JsonConvert.DeserializeObject(fortuneJson);

                string code = json["code"].ToString();
                string msg = json["msg"].ToString();
                string info = json["info"].ToString();
                string warn = json["warn"].ToString();

                string message = "";
                string picPath = Path.Combine(XQAPI.AppDir, "output.jpg");

                if (code != "200")
                {
                    if (msg != "")
                    {
                        message += msg;
                        e.FromGroup.SendMessage(e.RobotQQ, message);
                        return;
                    } else
                    {
                        message += "[fortune-运势] 服务器失联中......";
                        e.FromGroup.SendMessage(e.RobotQQ, message);
                        return;
                    }
                }
                if (msg != "success")
                {
                    message += msg;
                    e.FromGroup.SendMessage(e.RobotQQ, message);
                    return;
                }
                if (warn != "")
                {
                    pic(apipic, fromDataParm, headerParm);
                    message += warn;
                    message += "[pic=" + picPath + "]";
                    e.FromGroup.SendMessage(e.RobotQQ, message);
                    return;
                }
                if (info != "")
                {
                    pic(apipic, fromDataParm, headerParm);
                    if (notSend())
                    {
                        message += info;
                    }
                    message += "[pic=" + picPath + "]";
                    e.FromGroup.SendMessage(e.RobotQQ, message);
                    return;
                }
                if (code == "200")
                {
                    pic(apipic, fromDataParm, headerParm);
                    message += "[pic=" + picPath + "]";
                    e.FromGroup.SendMessage(e.RobotQQ, message);
                    return;
                }
            }
        }
        public static string get_autime()
        {
            TimeSpan ts = DateTime.Now.ToUniversalTime() - new DateTime(1970, 1, 1);
            string autime = Convert.ToString((int)ts.TotalSeconds);
            return autime;
        }

        public static string get_authkey(string aukey, string autime)
        {
            string aukeytime = aukey + "|" + autime;
            string authkey = "";
            MD5 md5 = MD5.Create();
            byte[] s = md5.ComputeHash(Encoding.UTF8.GetBytes(aukeytime));
            for (int i = 0; i < s.Length; i++)
            {
                authkey = authkey + s[i].ToString("x2");
            }
            return authkey;
        }

        public static string note = "";

        public static bool notSend()
        {
            string todayTime = DateTime.Now.ToShortDateString().ToString();
            if (note == "")
            {
                note = todayTime;
                return true;
            }
            else if (note != todayTime)
            {
                note = todayTime;
                return true;
            }
            return false;
        }

        public static string fortune(string api, FromDataParm fromDataParm, HeaderParm headerParm)
        {
            string postString = "";
            postString += string.Format("client={0}", fromDataParm.Client);
            postString += string.Format("&version={0}", fromDataParm.Version);
            postString += string.Format("&bot={0}", fromDataParm.Bot);
            postString += string.Format("&types={0}", fromDataParm.Types);
            postString += string.Format("&fromGroup={0}", fromDataParm.FromGroup);
            postString += string.Format("&fromQQ={0}", fromDataParm.FromQQ);
            postString += string.Format("&ask={0}", fromDataParm.Ask);
            postString += string.Format("&limit={0}", fromDataParm.Limit);

            byte[] postData = Encoding.UTF8.GetBytes(postString);

            WebClient webClient = new WebClient();

            webClient.Headers.Add("Content-Type", "application/x-www-form-urlencoded");
            webClient.Headers.Add("authkey", headerParm.Authkey);
            webClient.Headers.Add("autime", headerParm.Au_time);

            byte[] responseData = webClient.UploadData(api, "POST", postData);
            string path = Path.Combine(XQAPI.AppDir, "output.txt");
            File.WriteAllBytes(@path, responseData);

            string srcString = Encoding.UTF8.GetString(responseData);
            return srcString;
        }

        public static void pic(string api, FromDataParm fromDataParm, HeaderParm headerParm)
        {
            string postString = "";
            postString += string.Format("client={0}", fromDataParm.Client);
            postString += string.Format("&version={0}", fromDataParm.Version);
            postString += string.Format("&bot={0}", fromDataParm.Bot);
            postString += string.Format("&types={0}", fromDataParm.Types);
            postString += string.Format("&fromGroup={0}", fromDataParm.FromGroup);
            postString += string.Format("&fromQQ={0}", fromDataParm.FromQQ);
            postString += string.Format("&ask={0}", fromDataParm.Ask);
            postString += string.Format("&limit={0}", fromDataParm.Limit);

            byte[] postData = Encoding.UTF8.GetBytes(postString);

            WebClient webClient = new WebClient();

            webClient.Headers.Add("Content-Type", "application/x-www-form-urlencoded");
            webClient.Headers.Add("authkey", headerParm.Authkey);
            webClient.Headers.Add("autime", headerParm.Au_time);

            byte[] responseData = webClient.UploadData(api, "POST", postData);
            string path = Path.Combine(XQAPI.AppDir, "output.jpg");
            File.WriteAllBytes(@path, responseData);
        }
    }
}