package api

import (
	"database/sql"
	"fmt"
)

// ListPaginated updates receiver with a list of experiments using pagination.
func (experiments *Experiments) ListPaginated(page int, limit int) error {
	var experiment Experiment
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id limit %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}

// Find updates receiver with a record for a single experiment.
func (experiment *Experiment) Find(id string) error {
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.experiment_id = ?;"
	row := db.QueryRow(query, id)
	err = row.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(err)
		}
		return err
	}
	return nil
}

// DoseResponse updates receiver experiment with its corresponding dose/response data.
func (experiment *Experiment) DoseResponse() error {
	var doseResponse DoseResponse
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT dose, response FROM dose_responses WHERE experiment_id = ?;", experiment.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&doseResponse.Dose, &doseResponse.Response)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		experiment.DoseResponses = append(experiment.DoseResponses, doseResponse)
	}
	return nil
}

// CellDrugCombination updates receiver with a list of all experiments where a cell line and a drug have been tested.
func (experiments *Experiments) CellDrugCombination(cellID string, drugID string, typ string) error {
	var (
		cell Cell
		drug Drug
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	err = cell.Find(cellID, typ)
	if err != nil {
		return err
	}
	err = drug.Find(drugID, typ)
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, t.tissue_id, t.tissue_name, da.dataset_id, da.dataset_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.cell_id = ? AND e.drug_id = ?;"
	rows, _ := db.Query(query, cell.ID, drug.ID)
	defer rows.Close()
	for rows.Next() {
		var experiment Experiment
		err = rows.Scan(&experiment.ID, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		experiment.Cell.ID = cell.ID
		experiment.Cell.Name = cell.Name
		experiment.Drug = drug
		err = experiment.DoseResponse()
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}

// CellDatasetCombination updates receiver with a list of all experiments where a cell line and dataset have been tested.
func (experiments *Experiments) CellDatasetCombination(cellID string, datasetID string, typ string) error {
	var (
		cell    Cell
		dataset Dataset
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	err = cell.Find(cellID, typ)
	if err != nil {
		return err
	}
	err = dataset.Find(datasetID, typ)
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id WHERE e.cell_id = ? AND e.dataset_id = ?;"
	rows, _ := db.Query(query, cell.ID, dataset.ID)
	defer rows.Close()
	for rows.Next() {
		var experiment Experiment
		err = rows.Scan(&experiment.ID, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		experiment.Cell.ID = cell.ID
		experiment.Cell.Name = cell.Name
		experiment.Dataset = dataset
		err = experiment.DoseResponse()
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}