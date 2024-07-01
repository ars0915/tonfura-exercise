# 題目
	- 設計特定航空的訂機票的功能，規格：
		- 1. 可以按照起點、目的地、日期查詢班機狀態
		  2. 採取分頁方式返回可用班機清單、價格以及剩餘機位
		  3. 航空公司有超賣的慣例，功能也需要考量到超賣的情境
		  4. 設計表結構和索引、編寫主要程式碼、考慮大流量高並發情況(可以使用虛擬碼實現)。
	- 超賣補充：
		- 發生某一艙位超售無法安排所有乘客乘機，較高艙位有空餘座位時，航空公司一般會以免費升艙的方式安排。如較低艙位有空餘，則會徵集自願降艙旅客，給予經濟、里程補償的方式。發生機票超售、全艙未能安排所有乘客登機時，航空公司一般會徵集自願放棄座位的乘客，有時會以安排下一趟航班升艙、贈送里程或給予經濟補償的方式鼓勵乘客自願放棄座位。若無足夠乘客自願放棄座位，航空公司一般會按票價高低、級別高低順序安排座位，在艙位票價或會員級別等條件相同的情況下一般為先到先得，特殊旅客和後續有聯程航段的旅客也會被優先安排。
- # Spec
    - ## 定義
        - 班機：航班設定，需包含以下資訊
            - 起點
            - 目的地
            - 起飛時間
            - 艙等：
                - 類型：`頭等`、`商務`、`經濟`
                - 座位數
                - 價格
                - 可超賣數量
                - 狀態：`可購買`、`完售`
                    - `可購買`：賣出數量 < 總座位數 + 可超賣數量
                    - `完售`：賣出數量 >= 總座位數 + 可超賣數量
                - 賣出數
                - 出票數
        - 訂單：訂購後產生，不代表一定有座位，為了簡化系統先不處理座位安排
            - 班機編號
            - 乘客資訊
            - 艙等
            - 價格
            - 狀態：`待確認`、`已出票`
                - `待確認`：訂購後預設狀態，尚未保證登機
                - `已出票`：確認可登機
    - ## 情境
        - 1. 當使用者訂購機票時
            - a. 只能選擇狀態為`可購買`的班機及艙等，當選擇`完售`的艙等會回傳錯誤
            - b. 成功建立訂單，並更新該艙等的賣出數量
        - 2. 辦理出票時
            - a. 目前出票數 < 該艙等座位數 -> 更新訂單狀態為`已出票`
            - b. 目前出票數 >= 該艙等座位 -> 更新訂單狀態為`安排中`，進入超賣流程
        - 3. 超賣的可處理方式，此系統不考量里程、經濟補償，提供符合以下條件的選項
            - a. 當較高艙等有空位時，免費升艙
                - i. 較高艙等**未超賣**
                - ii. 臨近起飛時間，較高艙等的出票數 < 座位數
            - b. 當較低艙位有空位時，鼓勵降艙
                - i. 較低艙等**未超賣**
                - ii. 臨近起飛時間，較低艙等的出票數 < 座位數
            - c. 安排下一班
                - i. 優先選擇相同或較高艙等**未超賣**的座位
                - ii. 當 c.i. 沒有符合，選擇較低艙等**未超賣**的座位
                - iii. 當 c.ii 也無空位，安排至相同艙等
                - iv. 當相同艙等出票數 >= 該艙等座位數時，再根據 c. 順序檢查下一班是否符合
            - 流程
              [![](https://mermaid.ink/img/pako:eNqlVF1r01AY_ivhXEgHbclJctIuiFf1Yhd6ozdKQEKbrcW1GUmKzlCYVahBu83JyroNB45pp3QTp1hKhz_G5qNX-wue5iSnM8xt1OQib57zPM95z_uecyyQ1woqkMD8ovYkX1R0k7mfkysMfgw2kTBMjMzMECDPWu7WsbP-0W90_F9v_R8_R_s7brs-ar2ryZWQAy338wd378BZf-Otdd3DPX_wYvRly7fbXtcmcnfX9g77w9Pm2eCVf7LvNPruZo-5yTj9TxjE8dnAroV23AV2w9PVKexCQ94a9l4PeyvEbWzbtadNTvjL6_8WiuJe060y8BKvWCOD1d5BJ2YwqRH76Bpd_KdJaHGNzl2RhwETzs57p37iNBve922_U3daR0Q_lwv3pMGFnFG7OXaPc0IWH7Lc1ldCISnN5X6vPD9fn4mvcJmC5BNToEsUYXbxvPSozrTHtHFEHJZTjSrhftvAJGft2Dmy3dUNb_slKdtkpSyTSt3C9cefu0EEacRFkQ6j6THyIBjjKUugEaKReF4ZdZhq8XahTNz3C7iQuRHRAokBo-NNBjg6wEXHlCJ8dNgoIkRHhiIo2vhxlU4zU4NsJi9IgrKql5VSAV9-1pgtA7OollUZSDgsKPpjGciVGuYpVVO7t1zJA8nUq2oS6Fp1oQikeWXRwH_VpYJiqrmSsqAr5Rh6u1AyNZ2CavB7h1y4wb2bBEtKBUgWeAokKIhpgYNolmd5xIoil02CZSClMiibZgUOcRkIMyLihVoSPNM0PBeXRgKCWYHP8FmE0CwM7B4GY-NUa38AeS7Mqw?type=png)](https://mermaid.live/edit#pako:eNqlVF1r01AY_ivhXEgHbclJctIuiFf1Yhd6ozdKQEKbrcW1GUmKzlCYVahBu83JyroNB45pp3QTp1hKhz_G5qNX-wue5iSnM8xt1OQib57zPM95z_uecyyQ1woqkMD8ovYkX1R0k7mfkysMfgw2kTBMjMzMECDPWu7WsbP-0W90_F9v_R8_R_s7brs-ar2ryZWQAy338wd378BZf-Otdd3DPX_wYvRly7fbXtcmcnfX9g77w9Pm2eCVf7LvNPruZo-5yTj9TxjE8dnAroV23AV2w9PVKexCQ94a9l4PeyvEbWzbtadNTvjL6_8WiuJe060y8BKvWCOD1d5BJ2YwqRH76Bpd_KdJaHGNzl2RhwETzs57p37iNBve922_U3daR0Q_lwv3pMGFnFG7OXaPc0IWH7Lc1ldCISnN5X6vPD9fn4mvcJmC5BNToEsUYXbxvPSozrTHtHFEHJZTjSrhftvAJGft2Dmy3dUNb_slKdtkpSyTSt3C9cefu0EEacRFkQ6j6THyIBjjKUugEaKReF4ZdZhq8XahTNz3C7iQuRHRAokBo-NNBjg6wEXHlCJ8dNgoIkRHhiIo2vhxlU4zU4NsJi9IgrKql5VSAV9-1pgtA7OollUZSDgsKPpjGciVGuYpVVO7t1zJA8nUq2oS6Fp1oQikeWXRwH_VpYJiqrmSsqAr5Rh6u1AyNZ2CavB7h1y4wb2bBEtKBUgWeAokKIhpgYNolmd5xIoil02CZSClMiibZgUOcRkIMyLihVoSPNM0PBeXRgKCWYHP8FmE0CwM7B4GY-NUa38AeS7Mqw)  
        - 4. 放棄座位的訂單會以 3.c. 的順序被安排下一班
        - 5. 當有人放棄訂單時，**不主動**重新訂單，需經過乘客同意（此系統不實作通知功能
    - ## Known Issue
        - 臨近起飛時間時，高艙等旅客辦理出票可能會遇到因其他人升艙而沒座位的情況
            - 目前先以進入超賣流程處理
            - 可考慮是否以購買時的價格或艙等安排優先順序，可能會需要把已升艙旅客重新安排
- # 需實作
    - 考量系統複雜度先不實現使用者驗證及詳細乘客資訊
    - 班機清單 API
        - input
            - source
            - destination
            - departure_date
            - sort_by
            - page
            - limit
        - output
            - meta
                - record_count
                - page_count
                - absolute_page
                - page_size
            - data: array
                - id
                - source
                - destination
                - departure_time
                - classes
                  - id
                  - type
                  - seat_amount
                  - oversell_amount
                  - price
                  - status
    - 班機訂票 API
        - input
          - flight_id
          - class_id
          - price
          - user_id
          - amount
        - output
            - data: struct
                - id
                - flight_id
                - user_id
                - class_id
                - price
                - status
                - amount
    - 出票（包含超賣調整建議）
        - input
            - booking_id
        - output
            - data: struct
                - check_in_status (`success` or `failure`)
                - suggestion (may be null if `check_in_status` =  success)
                    - flight_id
                    - class_id
    - 放棄訂單 API (自動調整班機)
        - input
            - booking_id
        - output
            - data: struct
                - booking_id
                - flight_id
                - user_id
                - class_id
                - price
                - status
                - amount
    - 修改訂單 API
        - input
            - booking_id
            - flight_id  (optional)
            - class_type (optional)
            - status (optional)
        - output
            - data: struct
                - id
                - flight_id
                - user_id
                - class_id
                - price
                - status
                - amount
## Design
## DB schema
[![](https://mermaid.ink/img/pako:eNqdU02PnDAM_SuRz-yofAxsc257qSpV6q1CGnnBC9FAghJn1Skz_70BZroVg_qVE37P-D1b9giVqQkkkH2nsLHYl1qE99yppmUxLtH0vNIsVC0-f3zFHFulG-GMtxXdwTU5VhpZGf3KseopMANa9pYOyAt1Ka85VYfO_Ul3xhaLh0B9uLfEp4FW-Y6QD9gbr3nFmBeyjrpumx2squ5qma5eQVVL1fGg9F2R25AY2btVt0_GHCfy__tdtKehbTLekd0ktrr6F-dT4b9cD409rf6-rtf5LB4exHlcGpCiMppRabeRZMbbsKRo8YV-3ZbfZdwmHNhJ6jwbl-KJOqObIAQR9GR7VHU4grmfEril4Bhk-KzRHkso9SXkoWfz5aQrkGw9RWCNb1qQz9i5EPmhRqbrEa3Q97ViY3-CNIeflsOb7y-CATXIEb6BLOJdVhR5WiT7OMvfJhGcQKa7_E2WJFmS77N4n2XpJYLvxgShZBc_7vM4L9I4z9K8yB_nYl9ncjJ6-QF9tDBj?type=png)](https://mermaid.live/edit#pako:eNqdU02PnDAM_SuRz-yofAxsc257qSpV6q1CGnnBC9FAghJn1Skz_70BZroVg_qVE37P-D1b9giVqQkkkH2nsLHYl1qE99yppmUxLtH0vNIsVC0-f3zFHFulG-GMtxXdwTU5VhpZGf3KseopMANa9pYOyAt1Ka85VYfO_Ul3xhaLh0B9uLfEp4FW-Y6QD9gbr3nFmBeyjrpumx2squ5qma5eQVVL1fGg9F2R25AY2btVt0_GHCfy__tdtKehbTLekd0ktrr6F-dT4b9cD409rf6-rtf5LB4exHlcGpCiMppRabeRZMbbsKRo8YV-3ZbfZdwmHNhJ6jwbl-KJOqObIAQR9GR7VHU4grmfEril4Bhk-KzRHkso9SXkoWfz5aQrkGw9RWCNb1qQz9i5EPmhRqbrEa3Q97ViY3-CNIeflsOb7y-CATXIEb6BLOJdVhR5WiT7OMvfJhGcQKa7_E2WJFmS77N4n2XpJYLvxgShZBc_7vM4L9I4z9K8yB_nYl9ncjJ6-QF9tDBj)
## Architecture
```
		                                           
		+----------------+                         
		|                |      +------+           
		|   application  |+---->|  DB  |           
		|                |      +------+           
		+--------^-------+                         
		         |                                 
		         |mutex lock                       
		         |                                 
		         |                                 
		 +-------v------+                          
		 |              |                          
		 |redis-cluster |                          
		 |              |                          
		 +--------- ----+                          
		                                           
```

## Run
### 環境
```
ip=$(ipconfig getifaddr en0) docker-compose up -d --build
```
