namespace go product
namespace csharp Zen.DataAccess.Product

struct TShippingFee {
	1:required string warehouse;	//仓库
	2:required double fee;			//运费
	3:required double localFee;	//汇率转换后的运费
}

struct TProduct {
	1:optional i64 cid;						//淘宝商品分类id
	2:optional string vendorName;			//卖家名
	3:optional string productName;			//商品名
	4:optional double unitPrice;			//单价
	5:optional double shippingFee;			//运费
	6:optional string productUrl;			//商品url
	7:optional string productImage;			//商品图片
	8:optional string originCode;			//采购国家
	9:optional string shopName;				//店铺名
	10:optional string location;			//卖家地址	
	11:optional string aroundwWarehouse;	//附近仓库
	12:optional bool isShippingFee;			//是否免国内运费
	13:optional i32 favoritesItemId;		//收藏商品id
	14:optional i32 favoriteCatId;			//收藏分类id
	15:optional string specialHandlingFeeMessage;	//美国商品手续费描述
	16:optional double specialHandlingFeePercent;	//美国商品手续费
	17:optional list<string> propertyNames;			//商品规格描述
	18:optional list<TShippingFee> shippingFees;		//各仓库运费
	19:optional list<string> descriptionImages;		//商品详细图片
	20:optional list<TCharacteristicGroup> characteristicGroups; //可选sku列表
	21:optional list<TSku> skus;			//商品sku列表
	22:optional list<string> itemImgs;	
	23:optional bool isEZBuy;  //是否是特殊商品
	24:optional string priceSymbol; //货币符号
	25:optional string localUnitPrice;	//汇率转换后的价格
	26:optional double localShipmentFee;		//汇率转换后的运费
	27:optional string errMsg;				//错误信息
}

struct TProductExtension {
	1:optional i64 cid;						//淘宝商品分类id
	2:required string vendorName;			//卖家名
	3:required string productName;			//商品名
	4:required double unitPrice;			//单价
	5:required double shippingFee;			//运费
	6:required string productUrl;			//商品url
	7:required string productImage;			//商品图片
	8:required string originCode;			//采购国家
	9:required string shopName;				//店铺名
	10:required string location;			//卖家地址	
	11:required string aroundwWarehouse;	//附近仓库
	12:required bool isShippingFee;			//是否免国内运费
	13:required i32 favoritesItemId;		//收藏商品id
	14:required i32 favoriteCatId;			//收藏分类id
	15:required string specialHandlingFeeMessage;	//美国商品手续费描述
	16:required double specialHandlingFeePercent;	//美国商品手续费
	17:optional list<string> propertyNames;			//商品规格描述
	18:required list<TShippingFee> shippingFees;		//各仓库运费
	19:optional list<string> descriptionImages;		//商品详细图片
	20:optional list<TCharacteristicGroup> characteristicGroups; //可选sku列表
	21:optional list<TSku> skus;			//商品sku列表
	22:optional list<string> itemImgs;	
	23:optional bool isEZBuy;  //是否是特殊商品
	24:required string priceSymbol; //货币符号
	25:required string localUnitPrice;	//汇率转换后的价格
	26:required double localShipmentFee;		//汇率转换后的运费
	27:optional string errMsg;				//错误信息
	28:optional string eta;					//ETA时间范围
	29:optional bool displayShippingIcon;	//是否 显示 icon（特殊运输方式要显示icon）
	30:optional string altProductName;		//商品英文名称
	31:optional list<TCharacteristicGroup> altCharacteristicGroups; //可选英文sku列表
	32:required bool primeAvailable;		//是否支持prime模式购买
}


struct TSimpleProduct {
	1:required string productName;			//商品名
	2:required string productUrl;			//商品url
	3:required string productImage;			//商品图片
	4:required string originCode;			//采购国家
	5:required bool isEZBuy;  //是否是特殊商品
	6:required string localUnitPrice;	//汇率转换后的价格
}

struct TSku {
	1:required string price;		//价格
	2:required string properties;	//sku组合编号
	3:required string propertiesName;//sku名字
	4:required i64 quantity;		//数量
	5:required i64 skuId;			//sku Id
	6:required i64 skuSpecId;
	7:required string status;
	8:required i64 withHoldQuantity;
}

struct TCharacteristic {
	1:required string propkey;		//标识id
	2:required string actualValue;	//实际值
	3:required string remark;		//描述
	4:required string imageUrl;		//图片url
	5:required bool isSelected;		//是否选中
}

struct TCharacteristicGroup {
	1:required string name;		//标识id
	2:required list<TCharacteristic> characteristics;
}

struct TProductReviewDetail {
	1: required i32 id;					// 商品评论的id，默认为0
	2: required string productUrl;		// 商品url
	3: required i32 userId;				// 用户id
	4: required i32 rating;				// 满意度
	5: optional i32 helpfulCount;		// 此条评论的采纳数
	6: required string comment;			// 商品评论内容
	7: optional string pictures;		// `图片key`数组，即客户端可以自行计算出所需的任意规格的图片url
	8: required bool setHelpful;		//用户是否设置过helpful
	9: required string nickName;		//用户昵称
	10:required string headPortraits;	//用户头像
	11:optional string sku;				//商品sku
	12:required string createDate;		//创建日期
}

struct TProductReviewCount {
	1: required i32 all;		// 某个商品所有评论的个数
	2: required i32 hasPhoto;	// 某个商品含有图片评论的个数
}

struct SearchFilterField {
	1: required string name;
	2: required i32 productCount;
}

struct SearchFilter {
	1: required string name;
	2: required list<SearchFilterField> fields;
}

struct SearchFilterCond {
	1: required string filterName;
	2: required string fieldName;
}

struct SearchSortCond {
	1: required string sort;
	2: required bool isDesc;
}

struct SearchResult {
	1: list<TSimpleProduct> products;
	2: list<string> sorts;
	3: list<SearchFilter> filters;
}

exception TwitterUnavailable {
    1: string message;
}

service Product {
	TProduct GetProductDetail(1:string productUrl, 2:string purchaseSource) throws (1:TwitterUnavailable cond);
	oneway void Ping();
}
