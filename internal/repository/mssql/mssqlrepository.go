package mssqlstorage

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/encoding/unicode"
	"log"
	"remains_api/config"
	"remains_api/internal/domain"
	"strings"
	"time"
)

// var db *sql.DB

type Server struct {
	db *sqlx.DB
}

// InitDatabase - sets database connection configuration
func InitDatabase(config *config.Config) *Server {
	var err error
	connString := getConnString(config)

	log.Printf("Setting connection to db with configuration: %s \n", connString)

	server := &Server{}
	server.db, err = sqlx.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error opening connection: ", err.Error())
	}

	server.db.SetConnMaxLifetime(time.Minute * 4)
	server.db.SetMaxOpenConns(100)

	return server
}

// gets configuration and returns appropiate connection string
func getConnString(c *config.Config) string {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;",
		c.Server, c.Username, c.Password, c.Database)
	return connString
}

// will verify the connection is available or generate a new one
func (s *Server) getConnection() {

	err := s.db.Ping()
	if err != nil {
		log.Fatal("Could not ping db: ", err.Error())
	}
	log.Println("Ping successful")
}

func (s *Server) GetAll(userid string) ([]domain.Remains, error) {
	res := make([]domain.Remains, 0)
	query := `select 
					goods_name = g.NAME,
				    producer=p.name,
					country = c.NAME,
					mnn = s.NAME,
					price = l.PRICE_SAL,
					contractor = con.NAME,
					store = st.NAME,
					series = isnull(se.SERIES_NUMBER, 'Нет серии'),
					best_before = isnull(CONVERT(NVARCHAR,se.BEST_BEFORE,101),'нет с.г.'),
					remain = l.QUANTITY_REM
					from lot (nolock)l
					inner join goods g (nolock) on g.id_goods = l.ID_GOODS
					inner join substance (nolock)s on s.ID_SUBSTANCE = g.ID_SUBSTANCE
					inner join PRODUCER(nolock) p on p.ID_PRODUCER = g.ID_PRODUCER
					inner join country(nolock) c on c.ID_COUNTRY = p.ID_COUNTRY
					inner join store st(nolock) on st.id_store = l.id_store
					inner join CONTRACTOR(nolock) con on con.ID_CONTRACTOR = st.ID_CONTRACTOR
					left join series(nolock)se on se.ID_SERIES = l.ID_SERIES
					where l.QUANTITY_REM - l.QUANTITY_RES >0 and exists
					(select top 1 1 from user_2_contractor where id_user = @userId 
					                                         and contractor_id = st.id_contractor)`
	err := s.db.Select(&res, query, sql.Named("userId", userid))
	if err != nil {
		log.Println("failed to select: ", err.Error())
	}
	if len(res) == 0 {
		return nil, errors.New("No data")
	} else {
		return res, nil
	}

}
func (s *Server) GetFiltered(params domain.RemainRequest) ([]domain.Remains, error) {
	res := make([]domain.Remains, 0)
	query := `select 
					goods_name = g.NAME,
				    producer=p.name,
					country = c.NAME,
					mnn = s.NAME,
					price = l.PRICE_SAL,
					contractor = con.NAME,
					store = st.NAME,
					series = isnull(se.SERIES_NUMBER, 'Нет серии'),
					best_before = isnull(cast(se.BEST_BEFORE as varchar(10)),'нет с.г.'),
					remain = l.QUANTITY_REM
					from lot (nolock)l
					inner join goods g (nolock) on g.id_goods = l.ID_GOODS
					inner join substance (nolock)s on s.ID_SUBSTANCE = g.ID_SUBSTANCE
					inner join PRODUCER(nolock) p on p.ID_PRODUCER = g.ID_PRODUCER
					inner join country(nolock) c on c.ID_COUNTRY = p.ID_COUNTRY
					inner join store st(nolock) on st.id_store = l.id_store
					inner join CONTRACTOR(nolock) con on con.ID_CONTRACTOR = st.ID_CONTRACTOR
					left join series(nolock)se on se.ID_SERIES = l.ID_SERIES
					where l.QUANTITY_REM - l.QUANTITY_RES >0
					AND 1=1`
	searchQuery := buildSearchQuery(query, &params, "")
	//err := s.db.Select(&res, queryWhole)
	rows, err := s.db.NamedQuery(searchQuery, params)
	for rows.Next() {
		var resOne domain.Remains
		err = rows.StructScan(&resOne)
		if err != nil {
			return nil, err
		}
		res = append(res, resOne)
	}

	if err != nil {
		return nil, errors.New("smth wrong")
	}
	defer rows.Close()
	return res, nil
}

func buildSearchQuery(query string, searchParams *domain.RemainRequest, group string) string {
	sb := strings.Builder{}
	sb.WriteString(query)

	if searchParams.GoodsName != "" {
		sb.WriteString(`AND (g.name like '%'+:goods_name+'%')`)
	}
	if searchParams.MNN != "" {
		sb.WriteString(`AND (s.name like '%'+:mnn+'%'  )`)
	}
	if searchParams.Producer != "" {
		sb.WriteString(`AND ( p.name like '%'+:producer +'%')`)
	}

	if group != "" {
		sb.WriteString(`AND [group] = '` + group + `'`)
	}

	return sb.String()
}

func (s *Server) GetOnlyGroup(group string, params domain.RemainRequest) ([]domain.Remains, error) {
	res := make([]domain.Remains, 0)
	query := `select 
					goods_name = g.NAME,
				    producer=p.name,
					country = c.NAME,
					mnn = s.NAME,
					price = l.PRICE_SAL,
					contractor = con.NAME,
					store = st.NAME,
					series = isnull(se.SERIES_NUMBER, 'Нет серии'),
					best_before = isnull(cast(se.BEST_BEFORE as varchar(10)),'нет с.г.'),
					remain = l.QUANTITY_REM
					from lot (nolock)l
					inner join goods g (nolock) on g.id_goods = l.ID_GOODS
					inner join substance (nolock)s on s.ID_SUBSTANCE = g.ID_SUBSTANCE
					inner join PRODUCER(nolock) p on p.ID_PRODUCER = g.ID_PRODUCER
					inner join country(nolock) c on c.ID_COUNTRY = p.ID_COUNTRY
					inner join store st(nolock) on st.id_store = l.id_store
					inner join CONTRACTOR(nolock) con on con.ID_CONTRACTOR = st.ID_CONTRACTOR
					inner join CONTRACTOR_REMAIN_GROUP crg on crg.id_contractor = con.ID_CONTRACTOR
					left join series(nolock)se on se.ID_SERIES = l.ID_SERIES
					where l.QUANTITY_REM - l.QUANTITY_RES >0
					AND 1=1`
	searchQuery := buildSearchQuery(query, &params, group)
	//err := s.db.Select(&res, queryWhole)
	rows, err := s.db.NamedQuery(searchQuery, params)
	for rows.Next() {
		var resOne domain.Remains
		err = rows.StructScan(&resOne)
		if err != nil {
			return nil, err
		}
		res = append(res, resOne)
	}

	if err != nil {
		return nil, errors.New("smth wrong")
	}
	defer rows.Close()
	return res, nil
}

func (s Server) LoginUser(loginStruct domain.LoginStruct) (IDUser string, err error) {
	passwordhash := hashSum(loginStruct.Login, loginStruct.Password)
	userID := ""
	query := `SELECT TOP 1 id_user = isnull(cast(id_user as varchar(36)),'0')  FROM META_USER
	WHERE 1=1
	and password_hash = @password`
	e := s.db.QueryRow(query, sql.Named("password", passwordhash)).Scan(&userID)
	if e != nil {
		return "", errors.New("bad login:pas")
	} else {
		return userID, nil
	}
}

func hashSum(login string, password string) string {
	s := fmt.Sprint(login, password)
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	s1, _ := encoder.String(s)
	buff := []byte(s1)
	sum := md5.Sum(buff)
	hash := base64.StdEncoding.EncodeToString(sum[:])
	return hash
}
