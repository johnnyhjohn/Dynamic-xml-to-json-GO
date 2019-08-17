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
	Values []DynamicValues
}

type DynamicValues struct {
	Value map[string]interface{}
}

// Função para fazer append na lista de contents
func (contents *Contents) AddItemToDynamicList(dynamicValue DynamicValues) []DynamicValues {
	contents.Values = append(contents.Values, dynamicValue)
	return contents.Values
}

/**
*
*   @author Johnny John
*
*   @param { xmlquery.Node } - Estrutura xml que iremos ler e gravar em DynamicValues
*   @return {map[string]interface} - Retornamos um map de string que montamos dinamicamente
*
*   @desc
*   Método responsavel por enriquecer o objeto com os valores parseados do xml dinamicamente
 */
func (dynamicValue *DynamicValues) enrichDynamicValue(nodeToTransform *xmlquery.Node) map[string]interface{} {

	newValue := map[string]interface{}{}
	// Loop para percorrermos o node de PRODUCT
	for product := nodeToTransform.FirstChild; product != nil; product = product.NextSibling {
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
		dynamicValue.Value = newValue
	}

	return dynamicValue.Value
}

func main() {

	valuesToInput := []DynamicValues{}
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
		dynamicValue := DynamicValues{}
		// Alimentamos nossa variavel dynamicValue com os dados dentro da estrutura do Node "PRODUCT" do indice atual
		// Passamos o node como parâmetro no método para alimentar nossa variavel com os valores
		dynamicValue.enrichDynamicValue(productNode)
		// Atribuimos o array de product no content
		contentWithValues.AddItemToDynamicList(dynamicValue)
	}

	// Transformamos nosso array em json
	encjson, _ := json.Marshal(contentWithValues)
	// Imprimimos como umas string
	fmt.Println(string(encjson))
	fmt.Println("Successfully parsed products.xml")

	defer xmlFile.Close()
}
