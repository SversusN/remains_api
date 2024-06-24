package mssqlstorage

import (
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
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

func (s *Server) GetAll() ([]domain.Remains, error) {
	res := make([]domain.Remains, 0)
	query := `select 
					goods_name = g.NAME,
				    producer=p.name,
					country = c.NAME,
					mnn = s.NAME,
					price = l.PRICE_SAL,
					contractor = con.NAME,
					store = s.NAME,
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
					where l.QUANTITY_REM - l.QUANTITY_RES >0`
	err := s.db.Select(&res, query)
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
					store = s.NAME,
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
					store = s.NAME,
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
