package mixer

import (
	"errors"
	"log"
	"regexp"
)

// Tag определяет категорию строки
type Tag struct {
	//Name string
	Re *regexp.Regexp
}

// Ingredient определяет этап генерации
type Ingredient struct {
	//Name        string
	MustHave    []*Tag
	MustNotHave []*Tag
}

// IngredientMeta определяет обязательность ингридиента в рецепте
type IngredientMeta struct {
	Item        *Ingredient
	Probability float32
}

// Recipe определяет формулу формирования строк
type Recipe struct {
	//Name  string
	Steps []IngredientMeta
}

type RecipeBook struct {
	Tags        map[string]*Tag
	Ingredients map[string]*Ingredient
	Recipes     map[string]*Recipe
	RecipeNames []string
}

func InitBook() RecipeBook {
	rb := RecipeBook{}
	rb.initTags()
	rb.initIngredient()
	rb.initRecipes()
	return rb
}

func (book *RecipeBook) initTags() error {
	var tags = make(map[string]*Tag)

	tags["any"] = &Tag{regexp.MustCompile(`.`)}

	tags["startsWithAnd"] = &Tag{regexp.MustCompile(`^[Aa]nd`)}
	tags["startsWithBut"] = &Tag{regexp.MustCompile(`^[B]ut`)}
	tags["startsWithWhen"] = &Tag{regexp.MustCompile(`^[Ww]hen`)}
	//tags["startsWithWhy"] = &Tag{regexp.MustCompile(`^[Ww]hy`)}
	tags["startsWithYou"] = &Tag{regexp.MustCompile(`^[Yy]ou`)}
	tags["startsWithIf"] = &Tag{regexp.MustCompile(`^[Ii]f`)}
	tags["startsWithI"] = &Tag{regexp.MustCompile(`^I`)}

	tags["haveI"] = &Tag{regexp.MustCompile(`^I|\s[Ii]\s`)}
	tags["haveBecause"] = &Tag{regexp.MustCompile(`(?:[Bb]e)?[Cc]ause`)}
	tags["haveQuestion"] = &Tag{regexp.MustCompile(`\?$`)}

	tags["is2Words"] = &Tag{regexp.MustCompile(`^\S* \S*$`)}
	book.Tags = tags
	return nil
}

func (book *RecipeBook) initIngredient() error {
	var ingredients = make(map[string]*Ingredient)

	// Базовые ингридиенты = теги
	for key, tag := range book.Tags {
		ingredients[key] = &Ingredient{
			MustHave:    []*Tag{tag},
			MustNotHave: []*Tag{},
		}
	}

	book.Ingredients = ingredients
	return nil
}

func (book *RecipeBook) initRecipes() error {
	var recipes = make(map[string]*Recipe)

	recipes["iAndBecauseWhen"] = &Recipe{
		[]IngredientMeta{
			{book.mustProvideIngredient("startsWithI"), 1},
			{book.mustProvideIngredient("startsWithAnd"), 0.5},
			{book.mustProvideIngredient("startsWithWhen"), 1},
			{book.mustProvideIngredient("haveBecause"), 0.5},
		},
	}
	recipes["iBecause"] = &Recipe{
		[]IngredientMeta{
			{book.mustProvideIngredient("startsWithI"), 1},
			{book.mustProvideIngredient("haveBecause"), 1},
		},
	}
	recipes["ifI"] = &Recipe{
		[]IngredientMeta{
			{book.mustProvideIngredient("startsWithIf"), 1},
			{book.mustProvideIngredient("startsWithI"), 1},
		},
	}
	recipes["question"] = &Recipe{
		[]IngredientMeta{
			{book.mustProvideIngredient("haveQuestion"), 1},
			{book.mustProvideIngredient("any"), 1},
		},
	}
	recipes["any2-4Lines"] = &Recipe{
		[]IngredientMeta{
			{book.mustProvideIngredient("any"), 1},
			{book.mustProvideIngredient("any"), 1},
			{book.mustProvideIngredient("any"), 0.5},
			{book.mustProvideIngredient("any"), 0.5},
		},
	}
	recipes["anySingle"] = &Recipe{
		[]IngredientMeta{
			{book.mustProvideIngredient("any"), 1},
		},
	}

	// Построение списка рецептов
	names := make([]string, 0, 10)
	for key, _ := range recipes {
		names = append(names, key)
	}

	book.RecipeNames = names
	book.Recipes = recipes
	return nil
}

func (book *RecipeBook) mustProvideTag(name string) *Tag {
	item, exist := book.Tags[name]
	if exist {
		log.Println(book.Tags)
		panic(errors.New("Could not provide tag " + name))
	}
	return item
}

func (book *RecipeBook) mustProvideIngredient(name string) *Ingredient {
	item, exist := book.Ingredients[name]
	if !exist {
		log.Println(book.Ingredients)
		panic(errors.New("Could not provide ingredient " + name))
	}
	return item
}
