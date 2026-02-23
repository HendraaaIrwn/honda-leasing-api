package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"gorm.io/gen"
)

const (
	defaultQueryOutPath = "internal/domain/query"
)

func main() {
	outPath := flag.String("out", defaultQueryOutPath, "output path for generated query code")
	flag.Parse()

	projectRoot, err := findProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	resolvedOutPath, err := resolvePath(projectRoot, *outPath)
	if err != nil {
		log.Fatalf("invalid out path: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           resolvedOutPath,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.ApplyBasic(
		models.Province{},
		models.Kabupaten{},
		models.Kecamatan{},
		models.Kelurahan{},
		models.Location{},
		models.TemplateTask{},
		models.TemplateTaskAttribute{},
		models.OAuthProvider{},
		models.User{},
		models.UserOAuthProvider{},
		models.Role{},
		models.UserRole{},
		models.Permission{},
		models.RolePermission{},
		models.MotorType{},
		models.Motor{},
		models.MotorAsset{},
		models.Customer{},
		models.LeasingProduct{},
		models.LeasingContract{},
		models.LeasingTask{},
		models.LeasingTaskAttribute{},
		models.LeasingContractDocument{},
		models.PaymentSchedule{},
		models.Payment{},
	)

	g.Execute()
	log.Printf("gorm/gen completed successfully. generated query at: %s", resolvedOutPath)
}

func findProjectRoot() (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine current file location")
	}

	dir := filepath.Dir(thisFile)
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if stat, err := os.Stat(goModPath); err == nil && !stat.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found from %s", filepath.Dir(thisFile))
}

func resolvePath(projectRoot, pathValue string) (string, error) {
	trimmed := strings.TrimSpace(pathValue)
	if trimmed == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if filepath.IsAbs(trimmed) {
		return filepath.Clean(trimmed), nil
	}

	return filepath.Join(projectRoot, filepath.FromSlash(trimmed)), nil
}
