using System.Collections.Generic;
using LitJson;
using StarMatrix;

namespace StardustChapter
{
    public class TestNetMgr : Singleton<TestNetMgr>
    {
        private NetWorkLoad mNetWorkLoad;
        public string Domain="http://192.168.12.57:9999";

        public void DoFakePost(string api, string msg,Action onSuccess, Action onFailure)
        {
            if (mNetWorkLoad == null)
                mNetWorkLoad = new NetWorkLoad();
        
            SortedDictionary<string, string> dic = new SortedDictionary<string, string>
            {
                {"action", "test"},
                {"msg",msg}
            };
            mNetWorkLoad.HttpSend($"{Domain}/{api}", dic, (request, response) =>
            {
                if (response != null && response.IsSuccess && !string.IsNullOrEmpty(response.DataAsText))
                {
                    var json = CommonTools.Base64Unzip(response.DataAsText);
                    JsonData jsonData = JsonMapper.ToObject(json);
                    if (jsonData["error_code"].ToInt() == 1 && jsonData["error_msg"].ToString() == "success") //success
                    {
                        onSuccess?.Invoke();
                    }
                    else
                    {
                        onFailure?.Invoke();
                    }
                };
            });
        }
    }   
}