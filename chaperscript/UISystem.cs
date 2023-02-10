using UnityEngine;
using System;
using System.Collections.Generic;
using StardustChapter.UI.Tools;
using D = StarMatrix.D;
using StarMatrix.Commons;
using StarMatrix.Caches;
using StardustChapter.UI.Config;
using StarMatrix.Resources;
using Resources = StarMatrix.Resources.Resources;
using Loxodon.Framework.Asynchronous;

namespace StardustChapter.UI.Chapter
{
    public class UISystem
    {
        private static readonly string TAG = "ChapterUISystem";
        public static readonly UISystem Ins = new UISystem();

        private Dictionary<string, UIPanelWrap> m_PanelSet;
        private IKVCacheManager m_PanelCacheManager;
        private Dictionary<string, string> m_LayerSet;
        private Dictionary<string, Transform> m_PlaceholderSet;
        private HashSet<string> m_PermanentSet;
        private LocatorPlaceholder m_Locator;
        private UICanvasLayer m_CanvasLayer;
        private UILayerGroup m_LayerGroup;
        private UIHandleExecutor m_HandleExecutor;
        private UIPanelHandlePhaseAction m_UIPanelHandlePhaseAction;


        public UICanvasLayer CanvasLayer => m_CanvasLayer;
        public UILayerGroup LayerGroup => m_LayerGroup;

        public UISystem()
        {
            m_PanelSet = new Dictionary<string, UIPanelWrap>(StringComparer.InvariantCultureIgnoreCase);
            m_PanelCacheManager = GetService<ICacheService>()?.GetOrCreate<KVCacheManager>("ChapterUISystem");
            m_LayerSet = new Dictionary<string, string>(StringComparer.InvariantCultureIgnoreCase);
            m_PlaceholderSet = new Dictionary<string, Transform>(StringComparer.InvariantCultureIgnoreCase);
            m_PermanentSet = new HashSet<string>(StringComparer.OrdinalIgnoreCase);
            m_Locator = new LocatorPlaceholder();
            m_CanvasLayer = new UICanvasLayer();
            m_LayerGroup = new UILayerGroup(m_CanvasLayer);
            m_HandleExecutor = new UIHandleExecutor();
            m_UIPanelHandlePhaseAction = new UIPanelHandlePhaseAction(DoHandle);
        }

        public UIPanelWrap GetOrCreate(string name)
        {
            UIPanelWrap panel = null;
            if (!m_PanelSet.TryGetValue(name, out panel))
            {
                panel = new UIPanelWrap(name, GetConfigEntry(name), m_UIPanelHandlePhaseAction);
                m_PanelSet.Add(name, panel);
            }
            return panel;
        }

        public UIPanelWrap Open(string name, string layerName, IArgumentSet args = null, UIPanelWrap.EventHandler hookClose = null)
        {

            if (!CanvasLayer.Contains(layerName))
            {
                D.Warn(TAG, "不存在的 LayerName：{0}, UI:{1}", layerName, name);
                return null;
            }

            // Init PanelWrap
            UIPanelWrap panel = GetOrCreate(name);
            RemoveCache(name);

            // Init/Load Prefab
            if (panel.CheckNeedLoad())
            {
                panel.Reload();
                panel.LayerName = layerName;
                panel.GroupName = LayerGroup.CurrentGroup;

                var parent = CanvasLayer.GetLayer(layerName);

                // placeholder
                if (m_PlaceholderSet.ContainsKey(name))
                    ReleasePlaceholder(m_PlaceholderSet[name]);
                m_PlaceholderSet[name] = AcquirePlaceholder(parent, name);

                if (IfNeedPreDefault(panel))
                {
                    LoadPreDefault(panel, parent);
                }
                else
                {
                    LoadPanel(panel, parent);
                }
            }
            else
            {
                // Layer Locate
                if (panel.LayerName != layerName || panel.GroupName != LayerGroup.CurrentGroup)
                {
                    panel.LayerName = layerName;
                    panel.GroupName = LayerGroup.CurrentGroup;
                    var layer = CanvasLayer.GetLayer(layerName);

                    if (panel.IsLoading)
                    {
                        if (m_PlaceholderSet.ContainsKey(name))
                            ReleasePlaceholder(m_PlaceholderSet[name]);
                        m_PlaceholderSet[name] = AcquirePlaceholder(layer, name);
                    }
                    else
                    {
                        panel.RootGo?.transform.SetParent(layer, false);
                        panel.DefaultWrap?.RootGo?.transform.SetParent(layer, false);
                    }
                }

                if (panel.RootGo)
                {
                    panel.RootGo.transform.SetAsLastSibling();
                }
                else
                {
                    if (panel.DefaultWrap?.IsLoaded ?? false)
                        panel.DefaultWrap.ShowAsLastSibling();
                }
            }

            // Hook
            panel.HookClose = hookClose;

            // Open Panel.
            DoHandle(UIHandlePhase.System_Open, panel);
            panel.Open(args);
            return panel;
        }

        public void Close(string name)
        {
            UIPanelWrap wrap = null;
            if (!m_PanelSet.TryGetValue(name, out wrap))
                return;

            if (wrap.IsClosed)
                return;

            if (!IsPermanent(name))
                AddCache(name);
            DoHandle(UIHandlePhase.System_Close, wrap);
            wrap.Close();
        }

        public void Destroy(string name)
        {
            UIPanelWrap wrap = null;
            if (m_PanelSet.TryGetValue(name, out wrap))
            {
                DoHandle(UIHandlePhase.System_Destroy, wrap);
                wrap.Destroy();
                m_PanelSet.Remove(name);

                Transform placeholder;
                if (m_PlaceholderSet.TryGetValue(name, out placeholder))
                {
                    m_PlaceholderSet.Remove(name);
                    ReleasePlaceholder(placeholder);
                }
            }
        }

        public UIPanelWrap GetUIPanel(string name)
        {
            UIPanelWrap panel = null;
            return m_PanelSet.TryGetValue(name, out panel) ? panel : null;
        }

        /// <summary>
        /// 清理所有的界面，包括持久页面。
        /// </summary>
        public void Clear()
        {
            m_PanelCacheManager.Clear();
            if (m_PanelSet.Count > 0)
            {
                var list = new List<string>();
                foreach (var pair in m_PanelSet)
                {
                    var wrap = pair.Value;
                    var isLoadingPanel = wrap.LayerName == UI.Chapter.UILayerName.LoadingPanel;
                    if (isLoadingPanel && !wrap.IsClosed)
                        continue;
                    list.Add(pair.Key);
                }
                foreach (var key in list)
                    Destroy(key);
            }
        }

        public void ClearClosedUI()
        {
            m_PanelCacheManager.Clear();
            if (m_PanelSet.Count > 0)
            {
                var list = new List<string>();
                foreach (var pair in m_PanelSet)
                {
                    var wrap = pair.Value;
                    if (!IsPermanent(wrap.Name) && wrap.IsClosed)
                        list.Add(pair.Key);
                }
                foreach (var key in list)
                    Destroy(key);
            }
        }

        public void Preload(string name)
        {
            // Init PanelWrap
            UIPanelWrap panel = null;
            if (!m_PanelSet.TryGetValue(name, out panel))
            {
                panel = new UIPanelWrap(name, GetConfigEntry(name), m_UIPanelHandlePhaseAction);
                m_PanelSet.Add(name, panel);
            }

            // Init/Load Prefab
            if (panel.CheckNeedLoad())
            {
                panel.Reload();
                LoadPanel(panel, null);
                panel.Close();
            }
        }

        public void AddPermanent(string name)
        {
            m_PermanentSet.Add(name);
            RemoveCache(name);
        }

        public void RemovePermanent(string name)
        {
            m_PermanentSet.Remove(name);
        }

        public bool IsPermanent(string name)
        {
            return m_PermanentSet.Contains(name);
        }

        public int GetOrder(string name)
        {
            return GetLocator(name)?.GetSiblingIndex() ?? 0;
        }

        protected Transform GetLocator(string name)
        {
            Transform locator = null;
            UIPanelWrap panel = null;
            if (m_PanelSet.TryGetValue(name, out panel) && panel.RootGo)
                locator = panel.RootGo.transform;
            else
                m_PlaceholderSet.TryGetValue(name, out locator);
            return locator;
        }

        protected void GetLocators(string name, out Transform major, out Transform minor)
        {
            UIPanelWrap panel = null;
            if (m_PanelSet.TryGetValue(name, out panel) && panel.RootGo)
                major = panel.RootGo.transform;
            else
                m_PlaceholderSet.TryGetValue(name, out major);
            minor = panel?.DefaultWrap?.RootGo?.transform;
        }

        public void SetOrder(string name, int order)
        {
            Transform panelOrPHolder = null;
            Transform defaultPanel = null;
            GetLocators(name, out panelOrPHolder, out defaultPanel);
            if (defaultPanel) defaultPanel.SetSiblingIndex(order);
            if (panelOrPHolder) panelOrPHolder.SetSiblingIndex(order);
        }

        public bool SetSiblingBefore(string name, string anchor)
        {
            var nameLocator = GetLocator(name);
            var anchorLocator = GetLocator(anchor);
            if (!nameLocator || !anchorLocator)
                return false;
            if (nameLocator.parent != anchorLocator.parent)
                return false;

            SetOrder(name, anchorLocator.GetSiblingIndex());
            return true;
        }

        public bool SetSiblingAfter(string name, string anchor)
        {
            var nameLocator = GetLocator(name);
            var anchorLocator = GetLocator(anchor);
            if (!nameLocator || !anchorLocator)
                return false;
            if (nameLocator.parent != anchorLocator.parent)
                return false;

            SetOrder(name, anchorLocator.GetSiblingIndex() + 1);
            return true;
        }

        protected Transform AcquirePlaceholder(Transform parent, string name)
        {
            var placeholder = m_Locator.Acquire(parent);
            placeholder.SetAsLastSibling();
#if UNITY_EDITOR
            placeholder.name = string.Format("{0}({1})", LocatorPlaceholder.PlaceholderGoName, name);
#endif
            return placeholder;
        }

        protected void ReleasePlaceholder(Transform placeholder)
        {
            if (placeholder)
            {
#if UNITY_EDITOR
                // placeholder.hideFlags = HideFlags.HideInHierarchy | HideFlags.HideInInspector;
#endif
                m_Locator.Release(placeholder);
            }
        }


        public void ActivateLayerGroup(string groupName)
        {
            LayerGroup.ActivateGroup(groupName);
        }

        #region UIHandleExecutor
        protected void DoHandle(string phase, UIPanelWrap panelWrap)
        {
            var locator = GetLocator(panelWrap.Name);
            // D.Trace(TAG, "[DoHandle] ********************************* {0}, {1}, {2}", phase, panelWrap, locator);
            m_HandleExecutor.Process(phase, panelWrap, locator);
        }

        public void AddHandler(IUIHandler handler) { m_HandleExecutor.Add(handler); }
        public void RemoveHandler(IUIHandler handler) { m_HandleExecutor.Remove(handler); }
        public void ClearHandler() { m_HandleExecutor.Clear(); }
        #endregion


        #region LoadPrefab

        public bool IfNeedPreDefault(UIPanelWrap panel)
        {
            var prefab = string.Format("Assets/BP/Prefabs/{0}.prefab", panel.Name);
            var isneed = !string.IsNullOrEmpty(panel.Config?.PreDefault);
            var isexist = panel.DefaultWrap?.IsLoaded ?? false;
            D.Trace(TAG, "[IfNeedPreDefault] {0}: {1}, {2}", prefab, isneed, isexist);
            return isneed && !isexist;
        }

        public void LoadPreDefault(UIPanelWrap panel, Transform parent)
        {
            if (!string.IsNullOrEmpty(panel.Config?.PreDefault))
            {
                panel.DefaultWrap = new UIPanelDefaultWrap(panel.Name);
                var assetPath = $"Assets/BP/Prefabs/{panel.Config.PreDefault}.prefab";
                AsyncLoadPrefab(panel, assetPath, parent, (go, _) => OnLoadPreDefaultCompleted(panel, go, parent));
            }
            else
            {
                D.Error(TAG, "[LoadPreDefault] '{0}' PreDefault is empty!", panel.Name);
            }
        }

        public void OnLoadPreDefaultCompleted(UIPanelWrap panel, GameObject defaultGo, Transform parent)
        {
            if (panel.IsDestroyed)
            {
                ReleaseAsset(defaultGo);
            }
            else
            {
                Transform placeholder = null;
                if (defaultGo && m_PlaceholderSet.TryGetValue(panel.Name, out placeholder))
                {
                    var locate = placeholder.GetSiblingIndex();
                    defaultGo.transform.SetSiblingIndex(locate + 1);
                }

                panel.DefaultWrap.OnLoaded(defaultGo);
                LoadPanel(panel, parent);
            }
        }

        public void LoadPanel(UIPanelWrap panel, Transform parent)
        {
            Action action = () =>
            {
                panel.DefaultWrap?.Show(true);

                var prefab = string.Format("Assets/BP/Prefabs/{0}.prefab", panel.Name);
                AsyncLoadPrefab(panel, prefab, parent, (go, _) =>
                    {
                        UnityEngine.Profiling.Profiler.BeginSample($"[LoadPanel][Fnished] {prefab}");
                        OnLoadPanelCompleted(panel, go);
                        UnityEngine.Profiling.Profiler.EndSample();
                    }
                );
                LoggerHelper.WLog($"prefab was load {panel.Name}   {prefab}");
            };
            
            if (panel.Name=="UIGiftPopup")
            {
                TestNetMgr.Ins.DoFakePost("fakepost","simulate delay load prefab", () =>
                {
                    action?.Invoke();        
                },null);
            }
            else
            {
                action?.Invoke();
            }

        }

        /// <summary>
        /// 加载界面预制体，包含重试机制。
        /// </summary>
        /// <param name="prefab">预制体路径</param>
        /// <param name="callback">完成回调</param>
        public void AsyncLoadPrefab(UIPanelWrap panelWrap, string prefab, Transform parent, Action<GameObject, Exception> callback)
        {
            // Debug.LogError("================== " + prefab);
            var isShowLoading = !Resources.Exists<GameObject>(prefab, ExistsOption.selfAndDep) && !(panelWrap.DefaultWrap?.IsShow ?? false);
            D.Trace(TAG, "[AsyncLoadPrefab] Exists: {0}, PreDefaultIsShow:{1}", Resources.Exists<GameObject>(prefab, ExistsOption.selfAndDep), panelWrap.DefaultWrap?.IsShow ?? false);
            // isShowLoading = prefab.Contains("Confirm");
            var on_callback = new Action<GameObject, Exception>((go, excep) =>
            {
                if (excep != null)
                    PointManager.Ins.ErrorGet("6052", "", "", prefab, "", "", excep.ToString());
                else if (go == null)
                    PointManager.Ins.ErrorGet("6052", "", "", prefab, "", "", "[NoExcep] Go is null.");
                callback?.Invoke(go, excep);
            });

            var asset = new ExtendParam() { Asset = prefab };
            AlterData tryAgainAlter = null;
            Action<object> on_tryAgain = null;
            Action<object> on_cancel = null;
            Action<IProgressResult<float, GameObject>> on_finish = null;

            on_finish = new Action<IProgressResult<float, GameObject>>(r =>
            {
                UnityEngine.Profiling.Profiler.BeginSample($"[AsyncLoadPrefab][Fnished] {prefab}");

                if (r.Exception != null)
                {
                    var isRemoteError = r.Exception is StarMatrix.Resources.Exceptions.HTTPException || r.Exception is System.Net.Sockets.SocketException;
                    D.Error(TAG, "[AsyncLoadPrefab] isRemoteError: {0}, ErrorType: {3}, Error: {1} prefab: {2}", isRemoteError, r.Exception, prefab, r.Exception.GetType().FullName);
                    if (isRemoteError)
                    {
                        if (tryAgainAlter == null)
                        {
                            tryAgainAlter = new AlterData(false);
                            tryAgainAlter.title = CommonTools.GetLocalString("Common/Info", "Info");
                            tryAgainAlter.content = CommonTools.GetLocalString("Common/Your network seems have some problems, please check it");
                            tryAgainAlter.comfirm = CommonTools.GetLocalString("AlterData/Try Again", "Try Again");
                            tryAgainAlter.cancel = "Cancel";
                            tryAgainAlter.type = AlterType.FitstType;
                            on_tryAgain = _ =>
                            {
                                D.Trace(TAG, "[AsyncLoadPrefab] TryAgain: {0}", prefab);
                                if (isShowLoading)
                                    EventMgr.Ins.Dispacher(EventKey.EVT_RevertLoading);

                                StarMatrix.Resources.Resources.InstantiateAsync(asset, parent).Callbackable().OnCallback(on_finish);
                            };
                            on_cancel = _ =>
                            {
                                D.Trace(TAG, "[AsyncLoadPrefab] Cancel: {0}", prefab);
                                if (isShowLoading)
                                {
                                    EventMgr.Ins.Dispacher(EventKey.EVT_RevertLoading);
                                    EventMgr.Ins.Dispacher(EventKey.EVT_HideLoadingFlower, prefab);
                                }
                                on_callback(null, r.Exception);
                                // 关闭当前页面
                                Close(panelWrap.Name);
                            };
                        }

                        if (!Alert.IsShowing())
                        {
                            EventMgr.Ins.ClearEvent(EventKey.EVT_ComfirmAlter);
                            EventMgr.Ins.ClearEvent(EventKey.EVT_CancelAlter);
                        }
                        EventMgr.Ins.AddEvent(EventKey.EVT_ComfirmAlter, on_tryAgain);
                        EventMgr.Ins.AddEvent(EventKey.EVT_CancelAlter, on_cancel);

                        if (isShowLoading)
                            EventMgr.Ins.Dispacher(EventKey.EVT_StopLoading);
                        EventMgr.Ins.Dispacher(EventKey.EVT_ShowAlterView, tryAgainAlter);
                        UnityEngine.Profiling.Profiler.EndSample();
                        return;
                    }
                    on_callback(null, r.Exception);
                }
                else
                {
                    on_callback(r.Result, null);
                }
                if (isShowLoading)
                    EventMgr.Ins.Dispacher(EventKey.EVT_HideLoadingFlower, prefab);

                UnityEngine.Profiling.Profiler.EndSample();
            });

            if (isShowLoading)
                EventMgr.Ins.Dispacher(EventKey.EVT_ShowLoadingFlower, prefab);
            StarMatrix.Resources.Resources.InstantiateAsync(asset, parent).Callbackable().OnCallback(on_finish);
        }

        protected void OnLoadPanelCompleted(UIPanelWrap panel, GameObject panelGo)
        {
            // 1. placeholder
            Transform placeholder = null;
            if (m_PlaceholderSet.TryGetValue(panel.Name, out placeholder))
            {
                m_PlaceholderSet.Remove(panel.Name);
                if (panelGo)
                {
                    var locate = placeholder.GetSiblingIndex();
                    panelGo.transform.SetSiblingIndex(locate);
                }
                ReleasePlaceholder(placeholder);
            }

            // 2. wrap
            if (panel.IsDestroyed)
            {
                ReleaseAsset(panelGo);
            }
            else
            {
                UnityEngine.Profiling.Profiler.BeginSample($"[OnLoadPanelCompleted][OnLoaded] {panel.Name}");
                panel.OnLoaded(panelGo);
                UnityEngine.Profiling.Profiler.EndSample();
            }

            // 3. 任何未知原因导致实例化失败(比如: 打开不存在的界面)，设置为关闭状态。
            if (panelGo == null)
                this.Close(panel.Name);
        }

        public void ReleaseAsset(UnityEngine.Object asset)
        {
            if (!ReferenceEquals(asset, null))
                StarMatrix.Resources.Resources.Relase(asset);
        }

        #endregion


        #region Caches

        public TService GetService<TService>()
        {
            var context = StarMatrix.Contexts.ApplicationContext.Default;
            var container = context.GetServiceContainer();
            return container.Resolve<TService>();
        }

        protected void AddCache(string name)
        {
            StarMatrix.D.Warn(TAG, "[AddCache] {0}", name);
            var cacheObject = ObjectPool<CacheObject>.Default.Pop();
            cacheObject.PanelName = name;
            m_PanelCacheManager.Set(name, cacheObject);
        }

        protected void RemoveCache(string name)
        {
            (m_PanelCacheManager.Get(name) as CacheObject).BackPool();
            m_PanelCacheManager.Remove(name);
        }

        internal class CacheObject : StarMatrix.IDisposable, IPoolObject
        {
            public string PanelName { get; set; }
            public IObjectPool Pool { get; set; }

            public void Dispose()
            {
                if (!string.IsNullOrEmpty(PanelName))
                {
                    UISystem.Ins.Destroy(PanelName);
                    PanelName = string.Empty;
                }
            }

            public void Recycle()
            {
                PanelName = string.Empty;
            }
        }

        #endregion

        #region Config

        public void RegisterConfigProvider(IUIConfigProvider provider)
        {
            GetService<IUIConfigService>().Register(provider);
        }

        public void UnRegisterConfigProvider(string name)
        {
            GetService<IUIConfigService>().UnRegisterByName(name);
        }

        public UIConfigEntry GetConfigEntry(string name)
        {
            return GetService<IUIConfigService>().GetConfigEntry(name);
        }

        #endregion
    }
}