local payinfotb = ShopManager.createPayInfo(false)
CS.StardustChapter.EventMgr.Ins:Dispacher(CS.StardustChapter.EventKey.EVT_QuickPayInChapter, json.encode(payinfotb))