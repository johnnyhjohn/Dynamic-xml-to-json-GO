/*
*   @author Johnny John
*   @desc
*   Gerador de JSON via XML dinamico utilizando xmlQuery
 */
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
)

type Contents struct {
	Values []DinamicValues
}

type DinamicValues struct {
	Value map[string]interface{}
}

// Função para fazer append na lista de contents
func (contents *Contents) AddItemToDinamicList(dinamicValue DinamicValues) []DinamicValues {
	contents.Values = append(contents.Values, dinamicValue)
	return contents.Values
}

func main() {

	valuesToInput := []DinamicValues{}
	contentWithValues := Contents{valuesToInput}

	// Abertura do xml
	xmlFile, err := os.Open("products.xml")

	// tratamento se não conseguir ler o arquivo xml
	if err != nil {
		fmt.Println("Erro ao ler xml")
	}

	// Gravamos a arvore do xml para a leitura
	doc, erro := xmlquery.Parse(xmlFile)

	// Tratamento caso de erro na geração da arvore de nodes
	if erro != nil {
		fmt.Println("Erro ao gerar arvore do xml")
	}

	// Buscamos a tag de products que nos retornará uma tree dos nodes dentro dele,
	// pode ser feito dinamico colocando em uma função como parâmetro
	nodeProducts := xmlquery.FindOne(doc, "//PRODUCTS")

	// Buscamos a tag de products que nos retornará uma tree dos nodes dentro dele
	for node := nodeProducts.FirstChild; node != nil; node = node.NextSibling {
		// Buscamos o Node de product, pode ser feito dinamico colocando em uma função como parâmetro
		productNode := xmlquery.FindOne(node, "//PRODUCT")
		// Validação se não encontrou, pode retornar vazio caso não encontre ou haja espaços
		if productNode == nil {
			continue
		}
		dinamicValue := DinamicValues{}
		newValue := map[string]interface{}{}
		// Loop para percorrermos o node de PRODUCT
		for product := productNode.FirstChild; product != nil; product = product.NextSibling {
			// Eliminamos os espaços em branco do node e valor,
			// pode ser que a leitura retorne apenas espaçoes em branco
			nodeIsBlank := strings.TrimSpace(product.Data) == ""
			valueIsBlank := strings.TrimSpace(product.InnerText()) == ""
			// Verificamos se não estão vazios e atribuimos ao nosso newValue
			if nodeIsBlank && valueIsBlank {
				continue
			}
			// Atribuimos dinamicamente o nome da key com o nome da tag
			newValue[product.Data] = product.InnerText()
			// Atribuimos direto o valor ao Value pois a cada loop que ele passar aqui irá ter um valor a mais
			dinamicValue.Value = newValue
		}
		// Atribuimos o array de product no content
		contentWithValues.AddItemToDinamicList(dinamicValue)
	}

	// Transformamos nosso array em json
	encjson, _ := json.Marshal(contentWithValues)
	// Imprimimos como umas string
	fmt.Println(string(encjson))
	fmt.Println("Successfully parsed products.xml")

	defer xmlFile.Close()
}
