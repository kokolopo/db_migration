package utils

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ForeignClientSQL = `SELECT kyc.ClientID, usr.UserID, plan.PlanID FROM client_tbl_kyc kyc JOIN web_tbl_user usr ON usr.KycID = kyc.KycID JOIN client_tbl_plan plan ON plan.ClientID = kyc.ClientID WHERE kyc.Email = ?`
var DELETERELATEDDATASQL1 = `
	DELETE FROM client_tbl_plan_portfolio_unit u WHERE u.PlanID = ?
	DELETE FROM client_tbl_plan_redemption rd WHERE rd.PlanID = ?
	DELETE FROM client_tbl_plan_portfolio pp WHERE pp.PlanID = ?
`
var DELETERELATEDDATASQL2 = `
	DELETE FROM web_tbl_transaction t WHERE t.ClientID = ?
	DELETE FROM client_tbl_plan plan WHERE plan.ClientID = ?
	DELETE FROM client_tbl_detail_bankaccount ba WHERE ba.ClientID = ?
	DELETE FROM client_tbl_detail_bankaccount_draft bad WHERE bad.UserID = ?
	DELETE FROM client_tbl_kyc_draft kycd WHERE kycd.UserID = ?
	DELETE FROM client_tbl_kyc kyc WHERE kyc.ClientID = ?
	DELETE FROM web_tbl_user u WHERE u.ID = ?
`
var DELETEALLRELATEDDATASQL = `
DELETE FROM client_tbl_plan_portfolio_unit; DELETE FROM client_tbl_plan_redemption; DELETE FROM client_tbl_plan_portfolio; DELETE FROM web_tbl_transaction ; DELETE FROM client_tbl_plan ; DELETE FROM client_tbl_detail_bankaccount ; DELETE FROM client_tbl_detail_bankaccount_draft ; DELETE FROM client_tbl_kyc_draft; DELETE FROM client_tbl_kyc;  DELETE FROM web_tbl_user; 
`

type ForeignClient struct {
	ClientID string `json:"ClientID"`
	UserID   int    `json:"UserID"`
	PlanID   int    `json:"PlanID"`
}

func ClientBankaccountRule(db1 *gorm.DB, tablename string, v map[string]any) {
	var dataTableDB1 []map[string]any

	// dapatkan clientID berdasarkan userID sekian
	if err := db1.Table(tablename).Where("UserID = ?", v["UserID"]).Find(&dataTableDB1).Error; err != nil {
		panic(err)
	}

	v["ClientID"] = dataTableDB1[0]["ClientID"]
	delete(v, "UserID")

}

func PlanPortfolioRule(v map[string]any) {
	t := time.Time{}
	if v["NextAutoDebetDate"] == t {
		v["NextAutoDebetDate"] = nil
	}
}

func TransactionRule(v map[string]any) {
	// perlu disesuaikan :
	// 5 = Redemtion -> 3
	// 21 = One Off -> 1
	// 22 = Auto Debet Reguler -> 2
	switch v["TrxTypeID"] {
	case 5:
		delete(v, "TrxTypeID")
		v["TrxTypeID"] = 3
	case 21:
		delete(v, "TrxTypeID")
		v["TrxTypeID"] = 1
	case 22:
		delete(v, "TrxTypeID")
		v["TrxTypeID"] = 2
	}
}

func DeleteReletedDataRule(tx, db2 *gorm.DB, email string) (bool, error) {
	if email != "" {
		//get userID, ClientID, PlanID base on email
		var cred []ForeignClient
		var sql = fmt.Sprintf("SELECT kyc.ClientID, usr.UserID, plan.PlanID FROM client_tbl_kyc kyc JOIN web_tbl_user usr ON usr.KycID = kyc.KycID JOIN client_tbl_plan plan ON plan.ClientID = kyc.ClientID WHERE kyc.Email = %s", email)

		err := db2.Raw(sql).Scan(&cred).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		// running sql delete
		for _, v := range cred {
			if err := db2.Exec("DELETE FROM client_tbl_plan_portfolio_unit u WHERE u.PlanID = ?", v.PlanID).Error; err != nil {
				tx.Rollback()
				return false, err
			}
			if err := db2.Exec("DELETE FROM client_tbl_plan_redemption rd WHERE rd.PlanID = ?", v.PlanID).Error; err != nil {
				tx.Rollback()
				return false, err
			}
			if err := db2.Exec("DELETE FROM client_tbl_plan_portfolio pp WHERE pp.PlanID = ?", v.PlanID).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

		if err := db2.Exec("DELETE FROM web_tbl_transaction t WHERE t.ClientID = ?", cred[0].ClientID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_plan plan WHERE plan.ClientID = ?", cred[0].ClientID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_detail_bankaccount ba WHERE ba.ClientID = ?", cred[0].ClientID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_detail_bankaccount_draft bad WHERE bad.UserID = ?", cred[0].UserID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_kyc_draft kycd WHERE kycd.UserID = ?", cred[0].UserID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_kyc kyc WHERE kyc.ClientID = ?", cred[0].ClientID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if err := db2.Exec("DELETE FROM web_tbl_user usr WHERE usr.UserID = ?", cred[0].UserID).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if email == "" {
		if err := db2.Exec("DELETE FROM client_tbl_plan_portfolio_unit").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_plan_redemption").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_plan_portfolio").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM web_tbl_transaction").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_plan").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_detail_bankaccount").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_detail_bankaccount_draft").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_kyc_draft").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM client_tbl_kyc").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
		if err := db2.Exec("DELETE FROM web_tbl_user").Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return false, err
		}
	}

	return true, nil
}

func MigrateRelatedData(email string, tx, db1, db2 *gorm.DB) (bool, error) {
	// ambil data dari tabel; web_tbl_user, client_tbl_kyc, client_tbl_kyc_draft,
	// client_tbl_detail_bankaccount_draft, client_tbl_detail_bankaccount, client_tbl_plan,
	// client_tbl_plan_portfolio, client_tbl_plan_redemption,
	// client_tbl_plan_portfolio_unit, web_tbl_transaction

	if email != "" {
		var (
			cred []ForeignClient

			userData           map[string]any
			kycData            map[string]any
			kcyDraftData       map[string]any
			bankAccDraftData   []map[string]any
			bankAccData        []map[string]any
			userPlanData       []map[string]any
			trnsactionData     []map[string]any
			planPortfolioData  []map[string]any
			planRedemptionData []map[string]any
			userUnitData       []map[string]any
		)

		var sql = fmt.Sprintf("SELECT kyc.ClientID, usr.UserID, plan.PlanID FROM client_tbl_kyc kyc JOIN web_tbl_user usr ON usr.KycID = kyc.KycID JOIN client_tbl_plan plan ON plan.ClientID = kyc.ClientID WHERE kyc.Email = %s", email)
		err := db1.Raw(sql).Scan(&cred).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}
		fmt.Println(cred)

		err = db1.Raw(fmt.Sprintf("SELECT * FROM web_tbl_user WHERE UserID = %d", cred[0].UserID)).Scan(&userData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT * FROM client_tbl_kyc_draft WHERE UserID = %d", cred[0].UserID)).Scan(&kcyDraftData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT * FROM client_tbl_kyc WHERE ClientID = '%s'", cred[0].ClientID)).Scan(&kycData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT * FROM client_tbl_detail_bankaccount_draft WHERE UserID = %d", cred[0].UserID)).Scan(&bankAccDraftData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT * FROM client_tbl_detail_bankaccount WHERE UserID = %d", cred[0].UserID)).Scan(&bankAccData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT * FROM client_tbl_plan WHERE ClientID = '%s'", cred[0].ClientID)).Scan(&userPlanData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT tr.* FROM client_tbl_plan_portfolio_unit u JOIN client_tbl_plan plan on u.PlanID = plan.PlanID JOIN web_tbl_transaction tr ON u.TrxID = tr.TrxID WHERE u.PositionDate = (SELECT MAX(u2.PositionDate) FROM client_tbl_plan_portfolio_unit u2 WHERE u2.PortfolioID = u.PortfolioID AND u2.PlanID = u.PlanID AND EXISTS (SELECT 1 FROM client_tbl_plan p WHERE p.PlanID = u2.PlanID AND plan.ClientID = '%s')) AND u.UnitBalance > 0 ORDER BY u.PositionDate DESC", cred[0].ClientID)).Scan(&trnsactionData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT pp.* FROM client_tbl_plan_portfolio_unit u JOIN client_tbl_plan plan on u.PlanID = plan.PlanID JOIN client_tbl_plan_portfolio pp ON u.PlanPortfolioID = pp.PlanPortfolioID WHERE u.PositionDate = (SELECT MAX(u2.PositionDate) FROM client_tbl_plan_portfolio_unit u2 WHERE u2.PortfolioID = u.PortfolioID AND u2.PlanID = u.PlanID AND EXISTS (SELECT 1 FROM client_tbl_plan p WHERE p.PlanID = u2.PlanID AND plan.ClientID = '%s')) AND u.PlanPortfolioID > 0 AND u.UnitBalance > 0 ORDER BY u.PositionDate DESC", cred[0].ClientID)).Scan(&planPortfolioData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT pr.* FROM client_tbl_plan_portfolio_unit u JOIN client_tbl_plan plan on u.PlanID = plan.PlanID JOIN client_tbl_plan_redemption pr ON u.PlanRedemptionID = pr.PlanRedemptionID WHERE u.PositionDate = (SELECT MAX(u2.PositionDate) FROM client_tbl_plan_portfolio_unit u2 WHERE u2.PortfolioID = u.PortfolioID AND u2.PlanID = u.PlanID AND EXISTS (SELECT 1 FROM client_tbl_plan p WHERE p.PlanID = u2.PlanID AND plan.ClientID = '%s')) AND u.PlanRedemptionID > 0 AND u.UnitBalance > 0 ORDER BY u.PositionDate DESC", cred[0].ClientID)).Scan(&planRedemptionData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = db1.Raw(fmt.Sprintf("SELECT * FROM vw_portfolio_unit_client WHERE ClientID = '%s' AND UnitBalance > 0", cred[0].ClientID)).Scan(&userUnitData).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		// fmt.Println(userUnitData)

		// web_tbl_user
		t := time.Time{}
		if userData["UserLockDate"] == t {
			userData["UserLockDate"] = nil
		}
		if userData["LastLogin"] == t {
			userData["LastLogin"] = nil
		}
		userData["UserPhone"] = strings.ReplaceAll(userData["UserPhone"].(string), " ", "")
		if err := db2.Table("web_tbl_user").Create(&userData).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		// client_tbl_kyc_draft
		if err := db2.Table("client_tbl_kyc_draft").Create(&kcyDraftData).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		// client_tbl_kyc
		delete(kycData, "UserID")
		if err := db2.Table("client_tbl_kyc").Create(&kycData).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		//client_tbl_detail_bankaccount
		for _, ba := range bankAccData {
			ba["ClientID"] = cred[0].ClientID
			delete(ba, "UserID")
			if err := db2.Table("client_tbl_detail_bankaccount").Create(&ba).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

		//client_tbl_detail_bankaccount_draft
		if err := db2.Table("client_tbl_detail_bankaccount_draft").Create(&bankAccDraftData).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		// client_tbl_plan
		if err := db2.Table("client_tbl_plan").Create(&userPlanData).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		// web_tbl_transaction
		for _, v := range trnsactionData {
			switch v["TrxTypeID"] {
			case 5:
				delete(v, "TrxTypeID")
				v["TrxTypeID"] = 3
			case 21:
				delete(v, "TrxTypeID")
				v["TrxTypeID"] = 1
			case 22:
				delete(v, "TrxTypeID")
				v["TrxTypeID"] = 2
			}
			if err := db2.Table("web_tbl_transaction").Create(&v).Error; err != nil {
				tx.Rollback()
				return false, err
			}

			if err := db2.Exec("UPDATE web_tbl_transaction SET TrxTypeID = 3 WHERE TrxTypeID = 5").Error; err != nil {
				tx.Rollback()
				return false, err
			}
			if err := db2.Exec("UPDATE web_tbl_transaction SET TrxTypeID = 1 WHERE TrxTypeID = 21").Error; err != nil {
				tx.Rollback()
				return false, err
			}
			if err := db2.Exec("UPDATE web_tbl_transaction SET TrxTypeID = 2 WHERE TrxTypeID = 22").Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

		//client_tbl_plan_portfolio
		for _, v := range planPortfolioData {
			t := time.Time{}
			if v["NextAutoDebetDate"] == t {
				v["NextAutoDebetDate"] = nil
			}
			if err := db2.Table("client_tbl_plan_portfolio").Create(&v).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

		//client_tbl_plan_redemption
		if planRedemptionData != nil {
			if err := db2.Table("client_tbl_plan_redemption").Create(&planRedemptionData).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

		//vw_portfolio_unit_client
		for _, v := range userUnitData {
			delete(v, "ClientID")
			if err := db2.Table("client_tbl_plan_portfolio_unit").Create(&v).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

	}

	return true, nil
}

// err = db1.Exec("SELECT * FROM client_tbl_plan_portfolio WHERE ClientID = ?", cred[0].ClientID).Scan(&kycData).Error
// 		if err != nil {
// 			tx.Rollback()
// 			return false, err
// 		}
// 		err = db1.Exec("SELECT * FROM client_tbl_plan_redemption WHERE ClientID = ?", cred[0].ClientID).Scan(&kycData).Error
// 		if err != nil {
// 			tx.Rollback()
// 			return false, err
// 		}
// 		err = db1.Exec("SELECT * FROM client_tbl_plan_portfolio_unit WHERE ClientID = ?", cred[0].ClientID).Scan(&kycData).Error
// 		if err != nil {
// 			tx.Rollback()
// 			return false, err
// 		}
