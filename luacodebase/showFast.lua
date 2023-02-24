local payinfotb = ShopManager.createPayInfo(true)
payrecord.isOptionUnlock = true
UIUtil.dispacher(CS.StardustChapter.EventKey.EVT_ShowUIQuickPayFrame, json.encode(payinfotb))