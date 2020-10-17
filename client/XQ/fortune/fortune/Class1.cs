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

[assembly: XQPlugin("fortune", "木理", "1.0.0", "运势")]

namespace fortune
{
    [Plugin]
    public class Fortune
    {
        public struct FromDataParm
        {
            public string Client;
            public string Version;
            public string Types;
            public string FromQQ;
            public string Ask;
            public string Limit;
        };

        public struct HeaderParm
        {
            public string Authkey;
            public string Au_time;
        }

        [EnableEvent]
        public static void enable(object sender, XQEventArgs e)
        {
            XQAPI.OutPutLog("fortune-运势初始化完毕");
        }

        [GroupMsgEvent]
        public static void onGroupMsg(object sender, XQAppGroupMsgEventArgs e)
        {
            if (e.Message.Text == "运势")
            {
                e.FromGroup.SendMessage(e.RobotQQ, "少女祈祷中......");

                string aukey = "test";
                string apifortune = "http://127.0.0.1:8000/fortune";
                string apipic = "http://127.0.0.1:8000/fortune.jpg";

                FromDataParm fromDataParm;
                fromDataParm.Client = "xq";
                fromDataParm.Version = "3";
                fromDataParm.Types = "李清歌";
                fromDataParm.FromQQ = "test";
                fromDataParm.Ask = "运势";
                fromDataParm.Limit = "on";

                string autime = get_autime();
                string authkey = get_authkey(aukey, autime);
                HeaderParm headerParm;
                headerParm.Authkey = authkey;
                headerParm.Au_time = autime;

                string fortuneJson = fortune(apifortune, fromDataParm, headerParm);

                e.FromGroup.SendMessage(e.RobotQQ, fortuneJson);
                JObject json = (JObject)JsonConvert.DeserializeObject(fortuneJson);
                string code = json["code"].ToString();

                e.FromGroup.SendMessage(e.RobotQQ, code);

                pic(apipic, fromDataParm, headerParm);

                e.FromGroup.SendMessage(e.RobotQQ, "完毕！");
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
                authkey = authkey + s[i].ToString("X").ToLower();
            }
            return authkey;
        }

        public static string fortune(string api, FromDataParm fromDataParm, HeaderParm headerParm)
        {
            string postString = "";
            postString += string.Format("client={0}", fromDataParm.Client);
            postString += string.Format("&version={0}", fromDataParm.Version);
            postString += string.Format("&types={0}", fromDataParm.Types);
            postString += string.Format("&fromQQ={0}", fromDataParm.FromQQ);
            postString += string.Format("&ask={0}", fromDataParm.Ask);
            postString += string.Format("&limit={0}", fromDataParm.Limit);

            byte[] postData = Encoding.UTF8.GetBytes(postString);

            WebClient webClient = new WebClient();

            webClient.Headers.Add("Content-Type", "application/x-www-form-urlencoded");
            webClient.Headers.Add("authkey", headerParm.Authkey);
            webClient.Headers.Add("autime", headerParm.Au_time);

            byte[] responseData = webClient.UploadData(api, "POST", postData);
            string srcString = Encoding.UTF8.GetString(responseData);

            return srcString;
        }
        public static void pic(string api, FromDataParm fromDataParm, HeaderParm headerParm)
        {
            string postString = "";
            postString += string.Format("client={0}", fromDataParm.Client);
            postString += string.Format("&version={0}", fromDataParm.Version);
            postString += string.Format("&types={0}", fromDataParm.Types);
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
            XQAPI.OutPutLog(path);
            if (File.Exists(path))
            {
                File.Delete(path);
            }
            File.Create(path);
            File.WriteAllBytes(path, responseData);
        }
    }
}