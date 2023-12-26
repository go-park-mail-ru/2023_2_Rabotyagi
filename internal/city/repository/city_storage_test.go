package repository_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/pashagolub/pgxmock/v3"
)

func TestGetFullCity(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface)
		expectedResponse        any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT "city".id,"city".name FROM public."city"`).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).
						AddRow(uint64(1), "Moscow"))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.City{
				{ID: 1, Name: "Moscow"},
			},
		},
		{
			name: "test empty",
			behaviorCategoryStorage: func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT "city".id,"city".name FROM public."city"`).
					WillReturnRows(pgxmock.NewRows([]string{}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.City(nil),
		},
		{
			name: "test more data",
			behaviorCategoryStorage: func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT "city".id,"city".name FROM public."city"`).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).
						AddRow(uint64(1), "Moscow").
						AddRow(uint64(2), "St. Petersburg").
						AddRow(uint64(3), "Podolsk"))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.City{
				{ID: 1, Name: "Moscow"},
				{ID: 2, Name: "St. Petersburg"},
				{ID: 3, Name: "Podolsk"},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewCityStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCategoryStorage(catStorage, mockPool)

			response, err := catStorage.GetFullCities(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualTest(response, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSearchCity(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                string
		behaviorCityStorage func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface)
		expectedResponse    any
		searchInput         string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCityStorage: func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT city.id, city.name FROM public."city"`).
					WithArgs("%mos%").
					WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).
						AddRow(uint64(1), "Moscow"))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.City{
				{ID: 1, Name: "Moscow"},
			},
			searchInput: "Mos",
		},
		{
			name: "test empty",
			behaviorCityStorage: func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT city.id, city.name FROM public."city"`).
					WithArgs("%mos%").
					WillReturnRows(pgxmock.NewRows([]string{}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.City(nil),
			searchInput:      "Mos",
		},
		{
			name: "test more data",
			behaviorCityStorage: func(m *repository.CityStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT city.id, city.name FROM public."city"`).
					WithArgs("%m%").
					WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).
						AddRow(uint64(1), "Moscow").
						AddRow(uint64(2), "Manchester").
						AddRow(uint64(3), "Miami"))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.City{
				{ID: 1, Name: "Moscow"},
				{ID: 2, Name: "Manchester"},
				{ID: 3, Name: "Miami"},
			},
			searchInput: "M",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			cityStorage, err := repository.NewCityStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCityStorage(cityStorage, mockPool)

			response, err := cityStorage.SearchCity(ctx, testCase.searchInput)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualTest(response, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
